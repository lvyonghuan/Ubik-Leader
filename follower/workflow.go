package follower

// PreparingFollower is a method for preparing the follower.
// Use follower's caller to send a prepare signal to the follower.
// It returns a channel for blocking until the follower is ready.
// And it returns an error channel for handling errors during the preparation process.
func (f *Follower) PreparingFollower(errChan chan error) chan struct{} {
	var ch = make(chan struct{}, 1)
	go func() {
		c, err := f.caller.GetFollower(f.UUID)
		if err != nil {
			errChan <- err
			return
		}

		err = c.PreparingFollower()
		if err != nil {
			errChan <- err
			return
		} else {
			f.log.Debug("Follower " + f.UUID + " is ready to running")
			ch <- struct{}{}
		}
	}()

	return ch
}

func (f *Follower) RunningFollower(errChan chan error) chan struct{} {
	var ch = make(chan struct{}, 1)
	go func() {
		c, err := f.caller.GetFollower(f.UUID)
		if err != nil {
			errChan <- err
			return
		}

		err = c.RunningFollower()
		if err != nil {
			errChan <- err
			return
		}
		f.log.Debug("Follower " + f.UUID + " is running")
	}()

	return ch
}
