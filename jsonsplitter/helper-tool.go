package jsonsplitter 


import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

//OSPathSeparator versi string dari os.PathSeparator
const OSPathSeparator = string(os.PathSeparator)

//appendFilePath add file/ folder ke path. esensi nya base path bisa di akhiri dengan path separator.
// fungsi ini memastikan kalau misal not end with path separator, akan di tambahkan path separator. kalau sudah, tinggal gabung string saja
// <strong>basePath: <strong> folder awal untuk di tambahkan file atau folder
func appendFilePath(basePath string, fileOrFolder string) string {
	if strings.HasSuffix(basePath, OSPathSeparator) {
		return basePath + fileOrFolder
	}
	return basePath + OSPathSeparator + fileOrFolder
}

//appendFilePaths versi ini dengan path berupa array. misal nested folder
func appendFilePaths(basePath string, fileOrFolders ...string) string {
	var rtvl string
	if !strings.HasSuffix(basePath, OSPathSeparator) {
		rtvl = basePath
	} else {
		rtvl = basePath[0 : len(basePath)-2]
	}
	for _, thePath := range fileOrFolders {
		rtvl = rtvl + OSPathSeparator + thePath
	}
	return rtvl

}

//makeDirectoryHelper helper membuat directory. kalau directory tidak ada. ini tidak menyertakan pembuatan kalau directory nested
func makeDirectoryHelper(destinationFolder string, loggerEntry *logrus.Entry) (err error) {
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		loggerEntry.Info(fmt.Sprintf("Folder to create  [%s] not exists, creating now", destinationFolder))
		errMkdir := os.Mkdir(destinationFolder, os.ModePerm)
		if errMkdir != nil {
			loggerEntry.Error(fmt.Sprintf("Fail to create  directory[%s], error ", destinationFolder), errMkdir)
			return errMkdir
		}
	}
	return nil
}

//forceDeleteFolder delete folder yang tidak kosong. paksa
func forceDeleteFolder(foldername string, loggerEntry *logrus.Entry) (err error) {
	return os.RemoveAll(foldername)
}
