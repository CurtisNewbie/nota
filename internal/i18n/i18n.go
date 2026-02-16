package i18n

import "sync"

// Language represents supported languages
type Language string

const (
	LanguageEnglish Language = "en"
	LanguageChinese Language = "zh"
)

// Translation holds all translatable strings
type Translation struct {
	Menu struct {
		Note          string
		File          string
		View          string
		Language      string
		NewNote       string
		Import        string
		Export        string
		MinimizedMode string
		English       string
		Chinese       string
		Delete        string
	}
	Dialog struct {
		NoNoteSelected      string
		PleaseSelectNote    string
		UnsavedChanges      string
		SaveBeforeClosing   string
		SaveBeforeSwitching string
		SaveBeforeNewNote   string
		ExportSuccessful    string
		NoteExported        string
		ImportSuccessful    string
		ImportedNotes       string
		DeleteNote          string
		SureDelete          string
		UnsavedNote         string
		CannotDelete        string
		TitleCannotBeEmpty  string
		DuplicateNote       string
		OverwriteDuplicate  string
		DuplicateNotes      string
		OverwriteDuplicates string
		NoNotes             string
		NoNotesToExport     string
		Saved               string
		UnsavedChangesText  string
		NoNotesAvailable    string
	}
	Editor struct {
		TitlePlaceholder   string
		ContentPlaceholder string
		PlaceholderSearch  string
		Save               string
		Exit               string
	}
	Status struct {
		Saved          string
		UnsavedChanges string
	}
	Database struct {
		Location string
	}
}

var (
	currentLanguage = LanguageEnglish
	translations    map[Language]*Translation
	translationsMu  sync.RWMutex
)

func init() {
	translations = map[Language]*Translation{
		LanguageEnglish: getEnglishTranslation(),
		LanguageChinese: getChineseTranslation(),
	}
}

// getEnglishTranslation returns English translations
func getEnglishTranslation() *Translation {
	t := &Translation{}
	t.Menu.Note = "Note"
	t.Menu.File = "File"
	t.Menu.View = "View"
	t.Menu.Language = "Language"
	t.Menu.NewNote = "New Note"
	t.Menu.Import = "Import"
	t.Menu.Export = "Export"
	t.Menu.MinimizedMode = "Minimized Mode"
	t.Menu.English = "English"
	t.Menu.Chinese = "中文"
	t.Menu.Delete = "Delete"

	t.Dialog.NoNoteSelected = "No Note Selected"
	t.Dialog.PleaseSelectNote = "Please select a note"
	t.Dialog.UnsavedChanges = "Unsaved Changes"
	t.Dialog.SaveBeforeClosing = "You have unsaved changes. Do you want to save them before closing?"
	t.Dialog.SaveBeforeSwitching = "You have unsaved changes. Do you want to save them before switching?"
	t.Dialog.SaveBeforeNewNote = "You have unsaved changes. Do you want to save them before creating a new note?"
	t.Dialog.ExportSuccessful = "Export Successful"
	t.Dialog.NoteExported = "Note exported successfully"
	t.Dialog.ImportSuccessful = "Import Successful"
	t.Dialog.ImportedNotes = "Successfully imported %d notes"
	t.Dialog.DeleteNote = "Delete Note"
	t.Dialog.SureDelete = "Are you sure you want to delete this note?"
	t.Dialog.UnsavedNote = "Unsaved Note"
	t.Dialog.CannotDelete = "This note has not been saved yet and cannot be deleted"
	t.Dialog.TitleCannotBeEmpty = "Title cannot be empty"
	t.Dialog.CannotDelete = "This note has not been saved yet and cannot be deleted"
	t.Dialog.DuplicateNote = "Duplicate Note"
	t.Dialog.OverwriteDuplicate = "If a note with the same ID exists, do you want to overwrite it?"
	t.Dialog.DuplicateNotes = "Duplicate Notes"
	t.Dialog.OverwriteDuplicates = "If notes with the same ID exist, do you want to overwrite them?"
	t.Dialog.NoNotes = "No Notes"
	t.Dialog.NoNotesToExport = "There are no notes to export"
	t.Dialog.Saved = "Saved"
	t.Dialog.UnsavedChangesText = "Unsaved changes"
	t.Dialog.NoNotesAvailable = "No notes available. Click 'New Note' to create one."

	t.Editor.TitlePlaceholder = "Note Title"
	t.Editor.ContentPlaceholder = "Note content..."
	t.Editor.PlaceholderSearch = "Search notes..."
	t.Editor.Save = "Save"
	t.Editor.Exit = "Exit"

	t.Status.Saved = "Saved"
	t.Status.UnsavedChanges = "Unsaved changes"

	t.Database.Location = "DB: %s"

	return t
}

// getChineseTranslation returns Chinese translations
func getChineseTranslation() *Translation {
	t := &Translation{}
	t.Menu.Note = "笔记"
	t.Menu.File = "文件"
	t.Menu.View = "视图"
	t.Menu.Language = "语言"
	t.Menu.NewNote = "新建笔记"
	t.Menu.Import = "导入"
	t.Menu.Export = "导出"
	t.Menu.MinimizedMode = "最小化模式"
	t.Menu.English = "English"
	t.Menu.Chinese = "中文"
	t.Menu.Delete = "删除"

	t.Dialog.NoNoteSelected = "未选择笔记"
	t.Dialog.PleaseSelectNote = "请选择一个笔记"
	t.Dialog.UnsavedChanges = "未保存的更改"
	t.Dialog.SaveBeforeClosing = "您有未保存的更改。要在关闭前保存吗？"
	t.Dialog.SaveBeforeSwitching = "您有未保存的更改。要在切换前保存吗？"
	t.Dialog.SaveBeforeNewNote = "您有未保存的更改。要在创建新笔记前保存吗？"
	t.Dialog.ExportSuccessful = "导出成功"
	t.Dialog.NoteExported = "笔记导出成功"
	t.Dialog.ImportSuccessful = "导入成功"
	t.Dialog.ImportedNotes = "成功导入 %d 条笔记"
	t.Dialog.DeleteNote = "删除笔记"
	t.Dialog.SureDelete = "确定要删除此笔记吗？"
	t.Dialog.UnsavedNote = "未保存的笔记"
	t.Dialog.CannotDelete = "此笔记尚未保存，无法删除"
	t.Dialog.TitleCannotBeEmpty = "标题不能为空"
	t.Dialog.CannotDelete = "此笔记尚未保存，无法删除"
	t.Dialog.DuplicateNote = "重复笔记"
	t.Dialog.OverwriteDuplicate = "如果存在相同ID的笔记，是否覆盖？"
	t.Dialog.DuplicateNotes = "重复笔记"
	t.Dialog.OverwriteDuplicates = "如果存在相同ID的笔记，是否覆盖？"
	t.Dialog.NoNotes = "无笔记"
	t.Dialog.NoNotesToExport = "没有可导出的笔记"
	t.Dialog.Saved = "已保存"
	t.Dialog.UnsavedChangesText = "未保存的更改"
	t.Dialog.NoNotesAvailable = "没有可用的笔记。点击'新建笔记'创建一个。"

	t.Editor.TitlePlaceholder = "笔记标题"
	t.Editor.ContentPlaceholder = "笔记内容..."
	t.Editor.PlaceholderSearch = "搜索笔记..."
	t.Editor.Save = "保存"
	t.Editor.Exit = "退出"

	t.Status.Saved = "已保存"
	t.Status.UnsavedChanges = "未保存的更改"

	t.Database.Location = "数据库: %s"

	return t
}

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	translationsMu.Lock()
	defer translationsMu.Unlock()
	currentLanguage = lang
}

// GetLanguage returns the current language
func GetLanguage() Language {
	translationsMu.RLock()
	defer translationsMu.RUnlock()
	return currentLanguage
}

// GetTranslation returns the translation for the current language
func GetTranslation() *Translation {
	translationsMu.RLock()
	defer translationsMu.RUnlock()
	return translations[currentLanguage]
}

// T is a shorthand for GetTranslation()
func T() *Translation {
	return GetTranslation()
}
