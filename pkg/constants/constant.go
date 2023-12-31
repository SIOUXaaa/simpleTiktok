/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package constants

// connection information
const (
	//MYSQL配置
	MySQLDefaultDSN = "douyin:douyin123@tcp(127.0.0.1:13306)/douyin?charset=utf8&parseTime=True&loc=Local"
)

// constants in the project
const (
	UserTableName      = "users"
	FollowsTableName   = "follows"
	VideosTableName    = "videos"
	MessageTableName   = "messages"
	FavoritesTableName = "likes"
	CommentTableName   = "comments"

	VideoFeedCount       = 30
	FavoriteActionType   = 1
	UnFavoriteActionType = 2

	MinioEndPoint        = "8.130.69.85:9000"
	MinioPath            = "http://8.130.69.85:9000/"
	MinioVideoBucketName = "video"
	MinioImgBucketName   = "image"
	MinioAccessKeyID     = "BQC4APkLaJcrMGelVuw2"
	MinioSecretKey       = "iCVlLZ6w3DPhf8Hi7J9K2t13mFEJudWMn564zr9U"

	TestSign       = "测试账号！ offer"
	TestAva        = "avatar/test1.jpg"
	TestBackground = "background/test1.png"
)
