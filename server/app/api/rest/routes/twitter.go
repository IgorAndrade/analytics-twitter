package route

import (
	"encoding/json"

	"github.com/IgorAndrade/analytics-twitter/server/internal/usecase"
	"github.com/labstack/echo/v4"
)

func query(c echo.Context, ctn GetterDI) error {
	var service usecase.Search
	if err := ctn.Fill(usecase.TWITTER, service); err != nil {
		return err
	}

	query := c.QueryParam("query")
	c.Logger().Debug(query)

	var m map[string]string
	if err := json.Unmarshal([]byte(query), &m); err != nil {
		return err
	}

	return service.Find(m)

}
