package routing

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/app/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRouter(prHandler *handlers.PullRequestHandler, userHandler *handlers.UserHandler, teamHandler *handlers.TeamHandler, statsHandler *handlers.StatsHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/team/add", func(w http.ResponseWriter, r *http.Request) {
		teamHandler.CreateTeamHandler(w, r)
	})
	router.HandleFunc("/team/get", func(w http.ResponseWriter, r *http.Request) {
		teamHandler.GetTeamHandler(w, r)
	})
	router.HandleFunc("/users/setIsActive", func(w http.ResponseWriter, r *http.Request) {
		userHandler.SetIsActiveHandler(w, r)
	})
	router.HandleFunc("/users/getReview", func(w http.ResponseWriter, r *http.Request) {
		userHandler.GetReviewHandler(w, r)
	})
	router.HandleFunc("/pullRequest/create", func(w http.ResponseWriter, r *http.Request) {
		prHandler.CreatePullRequestHandler(w, r)
	})
	router.HandleFunc("/pullRequest/merge", func(w http.ResponseWriter, r *http.Request) {
		prHandler.MergePullRequestHandler(w, r)
	})
	router.HandleFunc("/pullRequest/reassign", func(w http.ResponseWriter, r *http.Request) {
		prHandler.ReassignPullRequestHandler(w, r)
	})
	router.HandleFunc("/stats/users", func(w http.ResponseWriter, r *http.Request) {
		statsHandler.UserStatsHandler(w, r)
	})
	router.HandleFunc("/stats/prs", func(w http.ResponseWriter, r *http.Request) {
		statsHandler.PrStatsHandler(w, r)
	})
	router.HandleFunc("/stats/overall", func(w http.ResponseWriter, r *http.Request) {
		statsHandler.OverallStatsHandler(w, r)
	})

	return router
}
