package di

func (c *Container) LeaderboardService() *leaderboard_service.Service {
	if c.leaderboardService == nil {
		c.leaderboardService = leaderboard_service.New(
			c.LeaderboardRepository(),
			c.StatisticsRepository(),
			c.UserRepository(),
		)
	}
	return c.leaderboardService
}
