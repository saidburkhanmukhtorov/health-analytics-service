package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	mongodb "github.com/health-analytics-service/health-analytics-service/storage/mongo"
	"github.com/stretchr/testify/assert"
)

func TestHealthRecommendationRepo(t *testing.T) {
	db := createMongoDBConnection(t)
	healthRecommendationRepo := mongodb.NewHealthRecommendationRepo(db)

	t.Run("CreateHealthRecommendation", func(t *testing.T) {
		testRecommendation := &health.HealthRecommendation{
			UserId:             uuid.NewString(),
			RecommendationType: "Exercise",
			Description:        "Engage in at least 30 minutes of moderate-intensity exercise most days of the week.",
			Priority:           2,
		}
		createdID, err := healthRecommendationRepo.CreateHealthRecommendation(context.Background(), testRecommendation)

		assert.NoError(t, err, "CreateHealthRecommendation should not return an error")
		assert.NotEmpty(t, createdID, "Created health recommendation should have a valid ID")
	})

	t.Run("GetHealthRecommendation", func(t *testing.T) {
		// 1. Create a recommendation to retrieve
		testRecommendation := &health.HealthRecommendation{
			UserId:             uuid.NewString(),
			RecommendationType: "Diet",
			Description:        "Consume a balanced diet rich in fruits, vegetables, and whole grains.",
			Priority:           1,
		}
		createdID, err := healthRecommendationRepo.CreateHealthRecommendation(context.Background(), testRecommendation)
		assert.NoError(t, err, "Creating health recommendation for GetHealthRecommendation test failed")
		assert.NotEmpty(t, createdID, "Created health recommendation should have a valid ID")

		// 2. Get the recommendation
		retrievedRecommendation, err := healthRecommendationRepo.GetHealthRecommendation(context.Background(), createdID)

		assert.NoError(t, err, "GetHealthRecommendation should not return an error")
		assert.NotNil(t, retrievedRecommendation, "GetHealthRecommendation response should not be nil")
		assert.Equal(t, testRecommendation.UserId, retrievedRecommendation.UserId)
		assert.Equal(t, testRecommendation.RecommendationType, retrievedRecommendation.RecommendationType)
		assert.Equal(t, testRecommendation.Description, retrievedRecommendation.Description)
		assert.Equal(t, testRecommendation.Priority, retrievedRecommendation.Priority)
	})

	t.Run("UpdateHealthRecommendation", func(t *testing.T) {
		// 1. Create a recommendation to update
		testRecommendation := &health.HealthRecommendation{
			UserId:             uuid.NewString(),
			RecommendationType: "Sleep",
			Description:        "Aim for 7-9 hours of quality sleep per night.",
			Priority:           3,
		}
		createdID, err := healthRecommendationRepo.CreateHealthRecommendation(context.Background(), testRecommendation)
		assert.NoError(t, err, "Creating health recommendation for UpdateHealthRecommendation test failed")
		assert.NotEmpty(t, createdID, "Created health recommendation should have a valid ID")

		// 2. Update the recommendation
		updateRecommendation := &health.HealthRecommendation{
			Id:                 createdID,
			UserId:             uuid.NewString(),
			RecommendationType: "Updated Sleep",
			Description:        "Get at least 8 hours of sleep.",
			Priority:           1,
		}
		err = healthRecommendationRepo.UpdateHealthRecommendation(context.Background(), updateRecommendation)
		assert.NoError(t, err, "UpdateHealthRecommendation should not return an error")

		// 3. Retrieve the recommendation and verify the update
		retrievedRecommendation, err := healthRecommendationRepo.GetHealthRecommendation(context.Background(), createdID)
		assert.NoError(t, err, "GetHealthRecommendation after update should not return an error")
		assert.Equal(t, updateRecommendation.UserId, retrievedRecommendation.UserId, "UserId should be updated")
		assert.Equal(t, updateRecommendation.RecommendationType, retrievedRecommendation.RecommendationType, "RecommendationType should be updated")
		assert.Equal(t, updateRecommendation.Description, retrievedRecommendation.Description, "Description should be updated")
		assert.Equal(t, updateRecommendation.Priority, retrievedRecommendation.Priority, "Priority should be updated")
	})

	t.Run("DeleteHealthRecommendation", func(t *testing.T) {
		// 1. Create a recommendation to delete
		testRecommendation := &health.HealthRecommendation{
			UserId:             uuid.NewString(),
			RecommendationType: "Hydration",
			Description:        "Drink plenty of water throughout the day.",
			Priority:           2,
		}
		createdID, err := healthRecommendationRepo.CreateHealthRecommendation(context.Background(), testRecommendation)
		assert.NoError(t, err, "Creating health recommendation for DeleteHealthRecommendation test failed")
		assert.NotEmpty(t, createdID, "Created health recommendation should have a valid ID")

		// 2. Delete the recommendation
		err = healthRecommendationRepo.DeleteHealthRecommendation(context.Background(), createdID)
		assert.NoError(t, err, "DeleteHealthRecommendation should not return an error")

		// 3. Attempt to retrieve the deleted recommendation (should fail)
		retrievedRecommendation, err := healthRecommendationRepo.GetHealthRecommendation(context.Background(), createdID)
		assert.Error(t, err, "GetHealthRecommendation after delete should return an error")
		assert.Nil(t, retrievedRecommendation, "GetHealthRecommendation response should be nil after delete")
	})

	t.Run("ListHealthRecommendations", func(t *testing.T) {
		// 1. Create some recommendations for a specific user
		userID := uuid.NewString()
		testRecommendations := []*health.HealthRecommendation{
			{
				UserId:             userID,
				RecommendationType: "Exercise 1",
				Description:        "Engage in regular physical activity.",
				Priority:           1,
			},
			{
				UserId:             userID,
				RecommendationType: "Diet 2",
				Description:        "Maintain a healthy diet.",
				Priority:           2,
			},
		}

		for _, recommendation := range testRecommendations {
			_, err := healthRecommendationRepo.CreateHealthRecommendation(context.Background(), recommendation)
			assert.NoError(t, err, "Creating health recommendation for ListHealthRecommendations test failed")
		}

		// 2. Test listing all recommendations for the user
		req := &health.ListHealthRecommendationsRequest{
			UserId: userID,
		}
		retrievedRecommendations, err := healthRecommendationRepo.ListHealthRecommendations(context.Background(), req)
		assert.NoError(t, err, "ListHealthRecommendations should not return an error")
		assert.NotNil(t, retrievedRecommendations, "ListHealthRecommendations response should not be nil")
		assert.GreaterOrEqual(t, len(retrievedRecommendations), 2, "Should have at least two health recommendations for the user")

		// 3. Test filtering by RecommendationType
		req = &health.ListHealthRecommendationsRequest{
			UserId:             userID,
			RecommendationType: "Exercise 1",
		}
		retrievedRecommendations, err = healthRecommendationRepo.ListHealthRecommendations(context.Background(), req)
		assert.NoError(t, err, "ListHealthRecommendations should not return an error")
		assert.NotNil(t, retrievedRecommendations, "ListHealthRecommendations response should not be nil")
		assert.Equal(t, 1, len(retrievedRecommendations), "Should have one health recommendation matching the filter")
		assert.Equal(t, "Exercise 1", retrievedRecommendations[0].RecommendationType, "RecommendationType should match the filter")

		// 4. Test filtering by Priority
		req = &health.ListHealthRecommendationsRequest{
			UserId:   userID,
			Priority: 2,
		}
		retrievedRecommendations, err = healthRecommendationRepo.ListHealthRecommendations(context.Background(), req)
		assert.NoError(t, err, "ListHealthRecommendations should not return an error")
		assert.NotNil(t, retrievedRecommendations, "ListHealthRecommendations response should not be nil")
		assert.Equal(t, 1, len(retrievedRecommendations), "Should have one health recommendation matching the filter")
		assert.Equal(t, int32(2), retrievedRecommendations[0].Priority, "Priority should match the filter")
	})
}
