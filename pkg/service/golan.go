package service

import (
    "golan/pkg/core"
    "golang.org/x/sys/windows/svc"
    "log"
    "time"
)

type GolanService struct {
    isPaused chan bool
}

func (gSvc *GolanService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
    const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
    changes <- svc.Status{State: svc.StartPending}
    changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
    gSvc.isPaused = make(chan bool)
    go gSvc.workerLoop(5 * time.Second)

    // Service loop
loop:
    for {
        select {
        case c := <-r:
            switch c.Cmd {
            case svc.Interrogate:
                changes <- c.CurrentStatus
            case svc.Stop, svc.Shutdown:
                break loop
            case svc.Pause:
                changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
                gSvc.isPaused <- true
            case svc.Continue:
                changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
                gSvc.isPaused <- false
            default:
                log.Printf("unexpected control request #%d", c)
            }
        }
    }
    changes <- svc.Status{State: svc.StopPending}
    return
}

func (gSvc *GolanService) workerLoop(duration time.Duration) {
    paused := false
    targets := core.Init()
    for {
        select {
        case paused = <-gSvc.isPaused:
            break
        default:
            if !paused {
                err := core.Verify(targets)
                if err != nil {
                    log.Println(err)
                }
            }
            time.Sleep(duration)
        }
    }
}
