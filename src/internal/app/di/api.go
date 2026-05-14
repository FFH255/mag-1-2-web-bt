package di

func (c *Container) Router() *proto.Router {
	if c.router == nil {
		corsConfig := c.cfg.CORS
		swaggerConfig := c.cfg.Swagger

		router := proto.NewRouter()

		// CORS
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  corsConfig.AllowAllOrigins,
			AllowOrigins:     corsConfig.AllowedOrigins,
			AllowMethods:     corsConfig.AllowedMethods,
			AllowHeaders:     corsConfig.AllowedHeaders,
			AllowCredentials: corsConfig.AllowCredentials,
			ExposeHeaders:    corsConfig.ExposedHeaders,
		}))

		// RequestID
		router.Use(request_id_middleware.New(c.Logger(), c.UUIDGenerator()))

		// Middlewares
		router.Middleware(
			c.AuthMiddleware(),
			c.RegistrationMiddleware(),
			c.RefreshTokenMiddleware(),
		)

		// Handlers
		router.Handle(
			swagger.New(swaggerConfig.Login, swaggerConfig.Password),
			c.UsersMeStatisticsPostHandler(),
			c.LeaderboardGetHandler(),
		)

		c.router = router
	}
	return c.router
}

func (c *Container) AuthMiddleware() proto.Middleware {
	if c.authMiddleware == nil {
		c.authMiddleware = auth_middleware.New(
			c.AuthService(),
			c.CookieManager(),
			c.Logger(),
		)
	}
	return c.authMiddleware
}

func (c *Container) RegistrationMiddleware() proto.Middleware {
	if c.registrationMiddleware == nil {
		c.registrationMiddleware = registration_middleware.New(
			c.AuthService(),
			c.CookieManager(),
		)
	}
	return c.registrationMiddleware
}

func (c *Container) RefreshTokenMiddleware() proto.Middleware {
	if c.refreshTokenMiddleware == nil {
		c.refreshTokenMiddleware = refresh_token_middleware.New(
			c.CookieManager(),
		)
	}

	return c.refreshTokenMiddleware
}

func (c *Container) UsersMeStatisticsPostHandler() *users_me_statistics_post_handler.Handler {
	if c.usersMeStatisticsPostHandler == nil {
		c.usersMeStatisticsPostHandler = users_me_statistics_post_handler.New(
			c.StatisticsService(),
			c.Logger(),
		)
	}
	return c.usersMeStatisticsPostHandler
}

func (c *Container) UsersMeLeaderboardPositionsGetHandler() *users_me_leaderboard_positions_get_handler.Handler {
	if c.usersMeLeaderboardPositionsGetHandler == nil {
		c.usersMeLeaderboardPositionsGetHandler = users_me_leaderboard_positions_get_handler.New(
			c.LeaderboardService(),
			c.Logger(),
		)
	}
	return c.usersMeLeaderboardPositionsGetHandler
}

func (c *Container) LeaderboardGetHandler() *leaderboard_get_handler.Handler {
	if c.leaderboardGetHandler == nil {
		c.leaderboardGetHandler = leaderboard_get_handler.New(
			c.LeaderboardService(),
			c.Logger(),
		)
	}
	return c.leaderboardGetHandler
}
