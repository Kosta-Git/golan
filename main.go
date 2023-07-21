package main

import (
    "golan/cmd"
    "golang.org/x/sys/windows"
    "log"
    "os"
    "strings"
    "syscall"
)

func main() {
    // Thanks to https://stackoverflow.com/a/59147866/6721372
    if !isAdmin() {
        askAdminRights()
    }

    if !isAdmin() {
        log.Fatalf("Please run as admin")
    }
    cmd.Execute()
}

func askAdminRights() {
    verb := "runas"
    exe, _ := os.Executable()
    cwd, _ := os.Getwd()
    args := strings.Join(os.Args[1:], " ")

    verbPtr, _ := syscall.UTF16PtrFromString(verb)
    exePtr, _ := syscall.UTF16PtrFromString(exe)
    cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
    argPtr, _ := syscall.UTF16PtrFromString(args)

    var showCmd int32 = 1 //SW_NORMAL

    err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
    if err != nil {
        log.Println(err)
    }
}

func isAdmin() bool {
    _, err := os.Open("\\\\.\\PHYSICALDRIVE0")
    return err == nil
}
