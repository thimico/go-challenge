// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// JOBSITY - Go(lang) restful starter kit
//
// API Docs for JOBSITY v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 2.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package api

import (
	"crypto/sha1"
	"os"

	"my-chat-jobsity-challenge/pkg/utl/zlog"

	"my-chat-jobsity-challenge/pkg/api/auth"
	al "my-chat-jobsity-challenge/pkg/api/auth/logging"
	at "my-chat-jobsity-challenge/pkg/api/auth/transport"
	"my-chat-jobsity-challenge/pkg/api/password"
	pl "my-chat-jobsity-challenge/pkg/api/password/logging"
	pt "my-chat-jobsity-challenge/pkg/api/password/transport"
	"my-chat-jobsity-challenge/pkg/api/user"
	ul "my-chat-jobsity-challenge/pkg/api/user/logging"
	ut "my-chat-jobsity-challenge/pkg/api/user/transport"

	"my-chat-jobsity-challenge/pkg/utl/config"
	"my-chat-jobsity-challenge/pkg/utl/jwt"
	authMw "my-chat-jobsity-challenge/pkg/utl/middleware/auth"
	"my-chat-jobsity-challenge/pkg/utl/postgres"
	"my-chat-jobsity-challenge/pkg/utl/rbac"
	"my-chat-jobsity-challenge/pkg/utl/secure"
	"my-chat-jobsity-challenge/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.New(os.Getenv("DATABASE_URL"), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.Service{}
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	authMiddleware := authMw.Middleware(jwt)

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec, rbac), log), e, authMiddleware)

	v1 := e.Group("/v1")
	v1.Use(authMiddleware)

	ut.NewHTTP(ul.New(user.Initialize(db, rbac, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, rbac, sec), log), v1)
	//ct.NewHTTP(chat.Initialize(db, rbac, sec)), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
