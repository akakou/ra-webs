package ta

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (ta *TA) Config(e *echo.Echo) error {
	output, err := ta.Register()
	if err != nil {
		return err
	}

	fmt.Printf("TTP: output %v\n", output)

	e.AutoTLSManager, err = ta.TLSConfig()
	if err != nil {
		return err
	}

	return nil
}
