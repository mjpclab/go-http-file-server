package i18n

type Translation struct {
	LoginLabel string
	MkdirLabel string

	UploadFilesLabel       string
	UploadFilesHint        string
	UploadDirLabel         string
	UploadDirHint          string
	UploadDirContentsLabel string
	UploadDirContentsHint  string
	UploadLabel            string
	UploadingLabel         string
	UploadSuccessLabel     string
	UploadFailLabel        string
	UploadDirFailMessage   string

	ListDirLabel  string
	ListNameLabel string
	ListTypeLabel string
	ListSizeLabel string
	ListTimeLabel string

	FilterLabel string

	SelectStart  string
	SelectCancel string
	SelectAll    string

	ArchiveLabel string

	DeleteLabel   string
	DeleteConfirm string

	Error401 string
	Error403 string
	Error404 string
	Error500 string
}

var translationEnUs = Translation{
	LoginLabel: "Login",
	MkdirLabel: "Create dir",

	UploadFilesLabel:       "Files",
	UploadFilesHint:        "Upload files",
	UploadDirLabel:         "Dir",
	UploadDirHint:          "Upload directory itself",
	UploadDirContentsLabel: "Dir contents",
	UploadDirContentsHint:  "Upload contents of directory",
	UploadLabel:            "Upload",
	UploadingLabel:         "Uploading...",
	UploadSuccessLabel:     "Upload success",
	UploadFailLabel:        "Upload failed",
	UploadDirFailMessage:   "Cannot upload directory. Make sure mkdir is enabled, and using a supported browser.",

	ListDirLabel:  "Dir",
	ListNameLabel: "Name",
	ListTypeLabel: "Type",
	ListSizeLabel: "Size",
	ListTimeLabel: "Time",

	FilterLabel: "filter...",

	SelectStart:  "Select",
	SelectCancel: "Cancel",
	SelectAll:    "Select all",

	ArchiveLabel: "Archive",

	DeleteLabel:   "Delete",
	DeleteConfirm: "Confirm delete?",

	Error401: "401 unauthorized",
	Error403: "403 resource is forbidden",
	Error404: "404 resource not found",
	Error500: "500 potential issue occurred",
}

var translationZhSimp = Translation{
	LoginLabel: "登录",
	MkdirLabel: "建目录",

	UploadFilesLabel:       "文件",
	UploadFilesHint:        "上传文件",
	UploadDirLabel:         "目录",
	UploadDirHint:          "上传目录自身",
	UploadDirContentsLabel: "目录内容",
	UploadDirContentsHint:  "上传目录下的内容",
	UploadLabel:            "上传",
	UploadingLabel:         "上传中……",
	UploadSuccessLabel:     "上传成功",
	UploadFailLabel:        "上传失败",
	UploadDirFailMessage:   "无法上传目录。 请确保mkdir已启用，并使用受支持的浏览器。",

	ListDirLabel:  "目录",
	ListNameLabel: "名称",
	ListTypeLabel: "类型",
	ListSizeLabel: "大小",
	ListTimeLabel: "时间",

	FilterLabel: "筛选……",

	SelectStart:  "选择",
	SelectCancel: "取消",
	SelectAll:    "全选",

	ArchiveLabel: "打包",

	DeleteLabel:   "删除",
	DeleteConfirm: "确认删除吗？",

	Error401: "401 未授权",
	Error403: "403 禁止访问资源",
	Error404: "404 资源不存在",
	Error500: "500 发生潜在错误",
}

var translationZhTrad = Translation{
	LoginLabel: "登入",
	MkdirLabel: "建目錄",

	UploadFilesLabel:       "檔案",
	UploadFilesHint:        "上傳檔案",
	UploadDirLabel:         "目錄",
	UploadDirHint:          "上傳目錄自身",
	UploadDirContentsLabel: "目錄內容",
	UploadDirContentsHint:  "上傳目錄下的內容",
	UploadLabel:            "上傳",
	UploadingLabel:         "上傳中……",
	UploadSuccessLabel:     "上傳成功",
	UploadFailLabel:        "上傳失敗",
	UploadDirFailMessage:   "無法上傳目錄。 請確保mkdir已啟用，並使用受支援的瀏覽器。",

	ListDirLabel:  "目錄",
	ListNameLabel: "名稱",
	ListTypeLabel: "類型",
	ListSizeLabel: "大小",
	ListTimeLabel: "時間",

	FilterLabel: "篩選……",

	SelectStart:  "選擇",
	SelectCancel: "取消",
	SelectAll:    "全選",

	ArchiveLabel: "打包",

	DeleteLabel:   "刪除",
	DeleteConfirm: "確認刪除嗎？",

	Error401: "401 未授權",
	Error403: "403 禁止訪問資源",
	Error404: "404 資源不存在",
	Error500: "500 發生潛在錯誤",
}
