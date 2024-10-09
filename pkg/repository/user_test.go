package repository_test

import (
	"testing"

	"taskmanagementapi/pkg/repository"
	"taskmanagementapi/pkg/utils/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCheckUserExistsByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("successfully found user", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "6705824f80a09eb0313f0e42"},
			{Key: "email", Value: "test@example.com"},
		}))
		ur := repository.NewUserRepository(mt.Client.Database("test"))
		result, err := ur.CheckUserExistsByEmail("test@example.com")
		assert.True(t, result)
		assert.NoError(t, err)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		ur := repository.NewUserRepository(mt.Client.Database("test"))

		result, err := ur.CheckUserExistsByEmail("nonexistent@example.com")

		assert.False(t, result)
		assert.NoError(t, err)
	})

	mt.Run("error during FindOne", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "some error",
		}))

		ur := repository.NewUserRepository(mt.Client.Database("test"))

		result, err := ur.CheckUserExistsByEmail("error@example.com")

		assert.False(t, result)
		assert.EqualError(t, err, "some error")
	})
}

func TestUserSignUp(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("successful user sign-up", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ur := repository.NewUserRepository(mt.Client.Database("test"))
		sampleUser := models.UserSignup{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}
		err := ur.UserSignUp(sampleUser)
		assert.NoError(t, err)
	})

	mt.Run("failure due to insertion error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "insertion error",
		}))

		ur := repository.NewUserRepository(mt.Client.Database("test"))
		sampleUser := models.UserSignup{
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		}
		err := ur.UserSignUp(sampleUser)
		assert.Error(t, err)
		assert.EqualError(t, err, "insertion error")
	})
}

func TestFindUserDetailsByEmail(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
    mt.Run("successfully found user by email", func(mt *mtest.T) {
        mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
            {Key: "_id", Value: "6705824f80a09eb0313f0e42"},
            {Key: "name", Value: "John Doe"},
            {Key: "email", Value: "johndoe@example.com"},
            {Key: "password", Value: "hashed_password"},
        }))

        ur := repository.NewUserRepository(mt.Client.Database("test"))
        userDetails, err := ur.FindUserDetailsByEmail("johndoe@example.com")
        expectedUser := models.UserDetails{
            ID:       "6705824f80a09eb0313f0e42",
            Name:     "John Doe",
            Email:    "johndoe@example.com",
            Password: "hashed_password",
        }
        assert.NoError(t, err)
        assert.Equal(t, expectedUser, userDetails)
    })

    mt.Run("user not found by email", func(mt *mtest.T) {
        mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))
        ur := repository.NewUserRepository(mt.Client.Database("test"))
        userDetails, err := ur.FindUserDetailsByEmail("nonexistent@example.com")
        assert.NoError(t, err)
        assert.Equal(t, models.UserDetails{}, userDetails)
    })

    mt.Run("error during FindOne operation", func(mt *mtest.T) {
        mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
            Code:    11000,
            Message: "some error",
        }))
        ur := repository.NewUserRepository(mt.Client.Database("test"))
        userDetails, err := ur.FindUserDetailsByEmail("error@example.com")
        assert.EqualError(t, err, "some error")
        assert.Equal(t, models.UserDetails{}, userDetails)
    })
}