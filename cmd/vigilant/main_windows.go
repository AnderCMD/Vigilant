// Command vigilant is a Windows tray app that keeps the system awake and your
// chat status "Active" by periodically simulating a harmless F15 key press.
package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"fyne.io/systray"

	"github.com/AnderCMD/Vigilant/assets"
	"github.com/AnderCMD/Vigilant/internal/autostart"
	"github.com/AnderCMD/Vigilant/internal/config"
	"github.com/AnderCMD/Vigilant/internal/vigil"
	"github.com/AnderCMD/Vigilant/internal/winapi"
)

// version is overridden at build time with -ldflags "-X main.version=...".
var version = "dev"

const repoURL = "https://github.com/AnderCMD/Vigilant"

var intervalChoices = []time.Duration{
	time.Minute, 3 * time.Minute, 5 * time.Minute, 10 * time.Minute, 15 * time.Minute,
}

type app struct {
	cfg     config.Config
	cfgPath string
	runner  *vigil.Runner
	cancel  context.CancelFunc
	// display mirrors cfg.KeepDisplayAwake for the pulse goroutine.
	display atomic.Bool
}

func main() {
	ok, err := winapi.AcquireSingleInstance()
	if err != nil {
		log.Printf("single-instance check failed: %v", err)
	}
	if !ok {
		return // another instance is already running
	}

	a := &app{}
	a.cfgPath, err = config.Path()
	if err != nil {
		log.Printf("cannot resolve config path: %v", err)
	}
	if a.cfg, err = config.Load(a.cfgPath); err != nil {
		log.Printf("using default settings: %v", err)
	}

	a.display.Store(a.cfg.KeepDisplayAwake)
	a.runner = vigil.New(a.pulse, a.cfg.Interval(), a.cfg.StartPaused)
	var ctx context.Context
	ctx, a.cancel = context.WithCancel(context.Background())
	go a.runner.Run(ctx)

	systray.Run(a.onReady, a.cancel)
}

// pulse is called by the runner on every tick.
func (a *app) pulse() error {
	// Both calls are attempted so one failing does not disable the other.
	err1 := winapi.PressF15()
	err2 := winapi.ResetIdleTimers(a.display.Load())
	if err1 != nil {
		return err1
	}
	return err2
}

func (a *app) save() {
	if a.cfgPath == "" {
		return
	}
	if err := config.Save(a.cfgPath, a.cfg); err != nil {
		log.Printf("cannot save settings: %v", err)
	}
}

func (a *app) refresh(toggle *systray.MenuItem) {
	if a.runner.Paused() {
		systray.SetIcon(assets.IconPaused)
		systray.SetTooltip("Vigilant - Paused")
		toggle.SetTitle("Resume")
	} else {
		systray.SetIcon(assets.IconActive)
		systray.SetTooltip(fmt.Sprintf("Vigilant - Active (every %s)", formatDuration(a.runner.Interval())))
		toggle.SetTitle("Pause")
	}
}

func (a *app) onReady() {
	systray.SetTitle("Vigilant")

	toggle := systray.AddMenuItem("Pause", "Pause or resume Vigilant")
	systray.AddSeparator()

	intervalMenu := systray.AddMenuItem("Interval", "Time between activity pulses")
	choices := make([]*systray.MenuItem, len(intervalChoices))
	for i, d := range intervalChoices {
		choices[i] = intervalMenu.AddSubMenuItemCheckbox(formatDuration(d), "", a.cfg.Interval() == d)
	}
	display := systray.AddMenuItemCheckbox("Keep display on", "Also prevent the screen from turning off", a.cfg.KeepDisplayAwake)
	startup := systray.AddMenuItemCheckbox("Start with Windows", "Launch Vigilant when you sign in", autostart.Enabled())
	systray.AddSeparator()

	about := systray.AddMenuItem("Vigilant "+version, "Open the project page")
	configItem := systray.AddMenuItem("Open settings folder", "")
	systray.AddSeparator()
	quit := systray.AddMenuItem("Exit", "Quit Vigilant")

	a.refresh(toggle)
	systray.SetOnTapped(func() { a.setPaused(!a.runner.Paused(), toggle) })

	// Fan every choice click into one channel so a single loop handles the menu.
	picked := make(chan int)
	for i, item := range choices {
		go func() {
			for range item.ClickedCh {
				picked <- i
			}
		}()
	}

	go func() {
		for {
			select {
			case <-toggle.ClickedCh:
				a.setPaused(!a.runner.Paused(), toggle)
			case i := <-picked:
				d := intervalChoices[i]
				a.cfg.IntervalSeconds = int(d / time.Second)
				a.runner.SetInterval(d)
				for j, c := range choices {
					if j == i {
						c.Check()
					} else {
						c.Uncheck()
					}
				}
				a.save()
				a.refresh(toggle)
			case <-display.ClickedCh:
				a.cfg.KeepDisplayAwake = !a.cfg.KeepDisplayAwake
				a.display.Store(a.cfg.KeepDisplayAwake)
				setChecked(display, a.cfg.KeepDisplayAwake)
				a.save()
			case <-startup.ClickedCh:
				want := !startup.Checked()
				if err := autostart.Set(want); err != nil {
					log.Printf("cannot change autostart: %v", err)
					continue
				}
				setChecked(startup, want)
			case <-about.ClickedCh:
				_ = winapi.Open(repoURL)
			case <-configItem.ClickedCh:
				if dir, err := config.Dir(); err == nil {
					_ = winapi.Open(dir)
				}
			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (a *app) setPaused(p bool, toggle *systray.MenuItem) {
	a.runner.SetPaused(p)
	a.refresh(toggle)
}

func setChecked(item *systray.MenuItem, on bool) {
	if on {
		item.Check()
	} else {
		item.Uncheck()
	}
}

func formatDuration(d time.Duration) string {
	if d%time.Minute != 0 {
		return d.String()
	}
	if m := int(d / time.Minute); m != 1 {
		return fmt.Sprintf("%d minutes", m)
	}
	return "1 minute"
}
