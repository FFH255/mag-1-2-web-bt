package di

type Container struct {
	logger        internal.Logger
	cfg           *config.Config
	postgres      *postgres.Database
	redis         *redis.Client
	uuidGenerator *uuid_generator.Generator
	// Repositories
	leaderboardRepository *leaderboard_repository.Repository
	// Services
	leaderboardService *leaderboard_service.Service
	// Handlers
	usersMeStatisticsPostHandler          *users_me_statistics_post_handler.Handler
	usersMeLeaderboardPositionsGetHandler *users_me_leaderboard_positions_get_handler.Handler
	leaderboardGetHandler                 *leaderboard_get_handler.Handler
	//Middleware
	authMiddleware         proto.Middleware
	registrationMiddleware proto.Middleware
	refreshTokenMiddleware proto.Middleware
	// Server
	router *proto.Router
	server *proto.Server
}

func (c *Container) Server() *proto.Server {
	if c.server == nil {
		cfg := c.cfg.Server

		c.server = proto.NewServer(
			c.Router(),
			cfg.Address,
			proto.MustUnmarshalDuration(cfg.ShutdownTimeout),
			proto.MustUnmarshalDuration(cfg.ReadTimeout),
			proto.MustUnmarshalDuration(cfg.WriteTimeout),
			proto.MustUnmarshalDuration(cfg.IdleTimeout),
		)
	}
	return c.server
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		cfg: cfg,
	}
}
