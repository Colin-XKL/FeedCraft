package craft

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/feeds"
)

// 通过文章原始内容的title 和content 字段, 计算哈希值, 来生成一个唯一的guid, 替换掉原feed中的可能有问题的guid

func GetGuidProcessor(transFunc TransFunc) FeedItemProcessor {
	return func(item *feeds.Item) error {
		guid, err := transFunc(item)
		if err != nil {
			return err
		}
		item.Id = guid
		return nil
	}
}

// 根据feed 中文章标题和内容生成md5作为guid, 如果几个字段都为空则返回随机值
func feedItemGuidGenerator(item *feeds.Item) (string, error) {
	if len(item.Title) == 0 && len(item.Content) == 0 && len(item.Description) == 0 {
		return uuid.New().String(), nil
	}

	combinedInput := title + content + description
	hash := getMD5Hash(combinedInput))
	return fmt.Sprintf("%x", hash), nil
}

func GetGuidCraftOptions() []CraftOption {
	craftOptions := []CraftOption{
		OptionTransformFeedItem(GetGuidProcessor(feedItemGuidGenerator)),
	}
	return craftOptions
}
