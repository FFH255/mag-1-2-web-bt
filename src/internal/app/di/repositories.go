package di

func (c *Container) LeaderboardRepository() *leaderboard_repository.Repository {
	if c.leaderboardRepository == nil {
		c.leaderboardRepository = leaderboard_repository.New(
			c.Redis().UnderlyingClient(),
		)
	}
	return c.leaderboardRepository
}
