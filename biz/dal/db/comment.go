package db

import (
	"gorm.io/gorm"
	"simpleTiktok/pkg/constants"
	"simpleTiktok/pkg/errno"
	"time"
)

type Comment struct {
	ID          int64          `json:"id"`
	UserId      int64          `json:"user_id"`
	VideoId     int64          `json:"video_id"`
	CommentText string         `json:"comment_text"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Comment) TableName() string {
	return constants.CommentTableName
}

func AddNewComment(comment *Comment) error {
	if ok, _ := CheckUserExitsById(comment.UserId); !ok {
		return errno.UserIsNotExistErr
	}
	if ok, _ := CheckVideoExitsById(comment.VideoId); !ok {
		return errno.VideoIsNotExistErr
	}
	err := DB.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckVideoExitsById(video_id int64) (bool, error) {
	var video Video
	if err := DB.Where("id = ?", video_id).Find(&video).Error; err != nil {
		return false, err
	}
	if video == (Video{}) {
		return false, nil
	}
	return true, nil
}

// DeleteCommentById delete comment by comment id
func DeleteCommentById(commentId int64) error {
	if ok, _ := CheckCommentExist(commentId); !ok {
		return errno.CommentIsNotExistErr
	}
	comment := &Comment{}
	err := DB.Where("id = ?", commentId).Delete(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckCommentExist(commentId int64) (bool, error) {
	comment := &Comment{}
	err := DB.Where("id = ?", commentId).Find(comment).Error
	if err != nil {
		return false, err
	}
	if comment.ID == 0 {
		return false, nil
	}
	return true, nil
}

func GetCommentListByVideoID(video_id int64) ([]*Comment, error) {
	var comment_list []*Comment
	if ok, _ := CheckVideoExistById(video_id); !ok {
		return comment_list, errno.VideoIsNotExistErr
	}
	err := DB.Table(constants.CommentTableName).Where("video_id = ?", video_id).Find(&comment_list).Error
	if err != nil {
		return comment_list, err
	}
	return comment_list, nil
}

func GetCommentCountByVideoID(video_id int64) (int64, error) {
	var count int64
	if err := DB.Model(&Comment{}).Where("video_id = ?", video_id).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
