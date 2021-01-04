package redis_dumper_test

import (
	redisdumper "github.com/BitMedia-IO/tool-collector/redis-dumper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v3"
	"testing"
)

func TestRedisItem_Insert(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6384",
		Password: "",
		DB:       0,
	})

	assert.Nil(t, client.Ping().Err())
	assert.Nil(t, client.FlushAll().Err())

	type testCase struct {
		item redisdumper.RedisItem
		want bool
	}

	testList := []interface{}{"test list item 1", "test list item 2", "test list item 3"}

	testCases := []testCase{
		{item: redisdumper.RedisItem{Key: "k:simple", Type: redisdumper.RedisImportTypeString, Value: "simple", TTL: 3333}, want: true},
		{item: redisdumper.RedisItem{Key: "k:empty", Type: redisdumper.RedisImportTypeString, Value: "", TTL: 3333}, want: true},
		{item: redisdumper.RedisItem{Key: "", Type: redisdumper.RedisImportTypeString, Value: "without key", TTL: 3333}, want: true},
		{item: redisdumper.RedisItem{Key: "h:simple", Type: redisdumper.RedisImportTypeHash, Value: map[string]interface{}{"test": "test string"}, TTL: 3333}, want: true},
		{item: redisdumper.RedisItem{Key: "l:simple", Type: redisdumper.RedisImportTypeList, Value: testList, TTL: 3333}, want: true},
	}
	for _, testCase := range testCases {
		t.Run("Desc: "+testCase.item.Key, func(t *testing.T) {
			assert.Nil(t, testCase.item.Insert(client))

			switch testCase.item.Type {
			case redisdumper.RedisImportTypeString:
				res, err := client.Get(testCase.item.Key).Result()
				assert.Nil(t, err)
				assert.Equal(t, testCase.item.Value, res)
			case redisdumper.RedisImportTypeHash:
				res, err := client.HGetAllMap(testCase.item.Key).Result()
				assert.Nil(t, err)
				assert.Equal(t, "test string", res["test"])
			case redisdumper.RedisImportTypeList:
				res, err := client.LRange(testCase.item.Key, 0, -1).Result()
				assert.Nil(t, err)
				assert.ElementsMatch(t, testList, res)
			}
		})
	}
}
