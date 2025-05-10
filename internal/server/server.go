package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	echoSession "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"

	db "dailycards/internal/database"
)

type Server struct {
	srv    *echo.Echo
	db     *db.Queries
	secret string
}

func New(q *db.Queries, secret string) *Server {
	return &Server{srv: echo.New(), db: q, secret: secret}
}

func (s *Server) Setup() {
	store := sessions.NewCookieStore([]byte(s.secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode, // используем Lax
		Secure:   false,                // **важно** на HTTP
	}
	s.srv.Use(echoSession.Middleware(store))
	
	s.srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderContentType},
		AllowCredentials: true,
	}))

	s.srv.Use(middleware.Logger(), middleware.Recover())

	s.srv.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "web",
	}))

	api := s.srv.Group("/api")
	api.POST("/users", s.CreateUser)
	api.POST("/login", s.HandleLogin)
	api.POST("/logout", s.HandleLogout)
	api.GET("/me", s.HandleMe)

	auth := api.Group("")
	auth.Use(s.SessionAuth)
	auth.POST("/packs", s.CreatePack)
	auth.GET("/packs", s.ListPacks)
	auth.DELETE("/packs/:id", s.DeletePack)
	auth.POST("/packs/:pack_id/cards", s.CreateCard)
	auth.GET("/packs/:pack_id/cards", s.ListCards)
	auth.GET( "/packs/:pack_id/repeat", s.RepeatPack)
	auth.POST("/packs/:pack_id/finish", s.FinishPack)
	auth.GET( "/stats",               s.UserStats)
	auth.GET("/user_stats",  s.UserStats)
	auth.DELETE("/packs/:pack_id/cards/:card_id", s.DeleteCard)
}

// Serve запускает HTTP‑сервер на :8080
func (s *Server) Serve() error {
	s.srv.Debug = true
	return s.srv.Start(":8080")
}

/* ------------------  USERS  ------------------ */

// CreateUserRequest описывает тело запроса на регистрацию
type CreateUserRequest struct {
    Username string `json:"username" form:"username"`
    Password string `json:"password" form:"password"`
}

// CreateUser создаёт нового пользователя и сразу инициализирует для него запись в user_stats
func (s *Server) CreateUser(c echo.Context) error {
    // 1) Считываем и валидируем входные данные
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "failed to read data: " + err.Error(),
        })
    }
    if req.Username == "" || req.Password == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "username and password must not be empty",
        })
    }

    // 2) Хэшируем пароль
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "hashing error: " + err.Error(),
        })
    }

    // 3) Создаём пользователя в БД
    user, err := s.db.CreateUser(c.Request().Context(), db.CreateUserParams{
        Username:     req.Username,
        PasswordHash: string(hashed),
    })
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return c.JSON(http.StatusConflict, map[string]string{
                "error": "user with this username already exists",
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "db error: " + err.Error(),
        })
    }

    // 4) Инициализируем статистику пользователя (таблица user_stats)
    //    SQL-генератор sqlc сгенерирует для этого метод CreateUserStats
    if err := s.db.CreateUserStats(c.Request().Context(), user.ID); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "failed to initialize user stats: " + err.Error(),
        })
    }

    // 5) Отдаём ответ об успешной регистрации
    return c.JSON(http.StatusCreated, map[string]string{
        "message": "user successfully created",
    })
}


/* ------------------  PACKS  ------------------ */

type CreatePackRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func (s *Server) CreatePack(c echo.Context) error {
    // 1) Разбор и валидация тела
    var req CreatePackRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "failed to read data: " + err.Error(),
        })
    }
    if req.Name == "" || req.Category == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "name and category must not be empty",
        })
    }

    // 2) Достаём user_id из сессии
    sess, _ := echoSession.Get("session", c)
    uidStr, ok := sess.Values["user_id"].(string)
    if !ok || uidStr == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }
    userID := uuidFromString(uidStr)

    // 3) Создаём Pack в БД
    pack, err := s.db.CreatePack(c.Request().Context(), db.CreatePackParams{
        Name:     req.Name,
        Category: pgtype.Text{String: req.Category, Valid: true},
    })
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return c.JSON(http.StatusConflict, map[string]string{
                "error": "pack with this name already exists",
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "db error: " + err.Error(),
        })
    }

    // 4) Увеличиваем счётчик созданных паков в user_stats
    if err := s.db.IncPacksCreated(c.Request().Context(), userID); err != nil {
        // не фатально, но логируем
        c.Logger().Warn("failed to increment packs_created:", err)
    }

    // 5) Отдаём клиенту ответ
    return c.JSON(http.StatusCreated, map[string]string{
        "message":  "pack successfully created",
        "id":       pack.ID.String(),
        "category": pack.Category.String,
    })
}

// ListPacks отдаёт все паки из БД
func (s *Server) ListPacks(c echo.Context) error {
    // (если нужны только свои паки) можно здесь посмотреть user_id из сессии
    packs, err := s.db.ListPacks(c.Request().Context())
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error: " + err.Error()})
    }
    return c.JSON(http.StatusOK, packs)
}


func (s *Server) DeletePack(c echo.Context) error {
    // 1) Получаем строку id из пути
    idParam := c.Param("id")
    if idParam == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "invalid pack id",
        })
    }

    // 2) Сканируем строку в pgtype.UUID (валидирует формат)
    var packID pgtype.UUID
    if err := packID.Scan(idParam); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "invalid pack id",
        })
    }

    // 3) Удаляем из БД
    if err := s.db.DeletePack(c.Request().Context(), packID); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "db error: " + err.Error(),
        })
    }

    // 4) Всё прошло успешно
    return c.NoContent(http.StatusNoContent)
}

/* ------------------  CARDS  ------------------ */

type CreateCardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Rating   *int32 `json:"rating,omitempty"`
}

func (s *Server) CreateCard(c echo.Context) error {
	packIDStr := c.Param("pack_id")
	var packID pgtype.UUID

	if err := packID.Scan(packIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid pack_id"})
	}

	var req CreateCardRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "cannot parse body: " + err.Error(),
		})
	}
	if req.Question == "" || req.Answer == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "question and answer are required",
		})
	}

	rating := pgtype.Int4{Int32: 0, Valid: true}
	if req.Rating != nil {
		rating = pgtype.Int4{Int32: *req.Rating, Valid: true}
	}

	card, err := s.db.CreateCard(c.Request().Context(), db.CreateCardParams{
		Question: req.Question,
		Answer:   req.Answer,
		PackID:   packID,
		Rating:   rating,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "pack not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, card)
}


func (s *Server) ListCards(c echo.Context) error {
	packIDParam := c.Param("pack_id")
	var packID pgtype.UUID

	if err := packID.Scan(packIDParam); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid pack_id",
		})
	}

	cards, err := s.db.ListCardsByPack(c.Request().Context(), packID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "db error: " + err.Error(),
		})
	}

	// 🧼 Чистый массив карт для фронта
	result := make([]map[string]interface{}, 0, len(cards))
	for _, card := range cards {
		result = append(result, map[string]interface{}{
			"id":       card.ID.String(),
			"question": card.Question,
			"answer":   card.Answer,
			"rating":   card.Rating.Int32, // Преобразуем pgtype.Int4 в обычное число
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (s *Server) DeleteCard(c echo.Context) error {
    // Берём user_id из сессии (чтобы не дать чужому удалить)
    sess, _ := echoSession.Get("session", c)
    uid, ok := sess.Values["user_id"].(string)
    if !ok || uid == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error":"unauthorized"})
    }

    // Получаем pack_id и card_id из пути (можно проверить оба, но далее удаляем только по card_id)
    packIDParam := c.Param("pack_id")
    cardIDParam := c.Param("card_id")

    var packID pgtype.UUID
    if err := packID.Scan(packIDParam); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error":"invalid pack_id"})
    }
    var cardID pgtype.UUID
    if err := cardID.Scan(cardIDParam); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error":"invalid card_id"})
    }

    // Удаляем карточку
    if err := s.db.DeleteCard(c.Request().Context(), cardID); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error":"db error: "+err.Error()})
    }

    return c.NoContent(http.StatusNoContent)
}

/* ------------------  AUTH  ------------------ */

func (s *Server) HandleLogin(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	user, err := s.db.GetUserByUsername(c.Request().Context(), req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	sess, _ := echoSession.Get("session", c)
	sess.Values["user_id"] = user.ID.String()

	// ★ Сохраняем сессию и игнорируем ошибку, чтобы не упало на 500
	if err := sess.Save(c.Request(), c.Response().Writer); err != nil {
		c.Logger().Warn("session save failed:", err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "logged in"})
}

func (s *Server) SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := echoSession.Get("session", c)
		if id, ok := sess.Values["user_id"].(string); !ok || id == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return next(c)
	}
}

func (s *Server) HandleMe(c echo.Context) error {
	sess, _ := echoSession.Get("session", c)
	idRaw, ok := sess.Values["user_id"].(string)
	if !ok || idRaw == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var uid pgtype.UUID
	if err := uid.Scan(idRaw); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid session"})
	}

	user, err := s.db.GetUserByID(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"username": user.Username})
}

func (s *Server) HandleLogout(c echo.Context) error {
	sess, _ := echoSession.Get("session", c)
	sess.Options.MaxAge = -1

	// ★ тоже сохраняем, чтобы браузер убрал куки
	_ = sess.Save(c.Request(), c.Response().Writer)
	return c.NoContent(http.StatusNoContent)
}

// ======================
// Интервальное повторение
// ======================

// RepeatPack возвращает карточки в нужном порядке
func (s *Server) RepeatPack(c echo.Context) error {
    var pid pgtype.UUID
    if err := pid.Scan(c.Param("pack_id")); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid pack_id"})
    }

    rows, err := s.db.ListRepeatCards(c.Request().Context(), pid)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error: " + err.Error()})
    }

    // ← здесь вместо return c.JSON(…, rows) делаем явное маппинг:
    out := make([]map[string]interface{}, 0, len(rows))
    for _, card := range rows {
        out = append(out, map[string]interface{}{
            "id":         card.ID.String(),
            "question":   card.Question,
            "answer":     card.Answer,
            "rating":     card.Rating.Int32,  // Int4 → int
            "last_wrong": card.LastWrong.Bool, // Bool → bool
        })
    }

    return c.JSON(http.StatusOK, out)
}

// FinishPack обрабатывает результаты сессии повторения
func (s *Server) FinishPack(c echo.Context) error {
    // 1) берём user_id из сессии
    sess, _ := echoSession.Get("session", c)
    uidStr, ok := sess.Values["user_id"].(string)
    if !ok || uidStr == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error":"unauthorized"})
    }

    // 2) читаем тело запроса
    var body struct {
        Stats []struct {
            CardID  string `json:"card_id"`
            Correct bool   `json:"correct"`
        } `json:"stats"`
    }
    if err := c.Bind(&body); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error":"invalid body"})
    }

    // 3) обновляем last_wrong, считаем дельту рейтинга
    delta := 0
    allCorrect := true
    for _, st := range body.Stats {
        // отметка о том, что в прошлый раз была ошибка
        _ = s.db.MarkCardWrong(c.Request().Context(), db.MarkCardWrongParams{
            LastWrong: pgtype.Bool{Bool: !st.Correct, Valid: true},
            ID:        uuidFromString(st.CardID),
        })

        if st.Correct {
            delta++
        } else {
            delta--
            allCorrect = false
        }
    }

    // 4) обновляем общий рейтинг пользователя
    _ = s.db.AddUserRating(c.Request().Context(), db.AddUserRatingParams{
        Rating: pgtype.Int4{Int32: int32(delta), Valid: true},
        UserID: uuidFromString(uidStr),
    })

    // 5) если ни разу не ошибился — +1 к mastered
    if allCorrect {
        _ = s.db.IncPacksMastered(c.Request().Context(), uuidFromString(uidStr))
    }

    return c.NoContent(http.StatusNoContent)
}

// UserStats отдаёт общую статистику по пользователю
func (s *Server) UserStats(c echo.Context) error {
    // 1) Достаём user_id из сессии
    sess, _ := echoSession.Get("session", c)
    uidRaw, ok := sess.Values["user_id"].(string)
    if !ok || uidRaw == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }

    // 2) Бежим в БД
    statsRow, err := s.db.GetUserStats(c.Request().Context(), uuidFromString(uidRaw))
    if err != nil {
        // если строки ещё нет — возвращаем нули и 200
        if errors.Is(err, pgx.ErrNoRows) {
            return c.JSON(http.StatusOK, map[string]interface{}{
                "rating":         0,
                "packs_created":  0,
                "packs_mastered": 0,
            })
        }
        // любая другая ошибка — 500
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "db error: " + err.Error(),
        })
    }

    // 3) Выдергиваем int из pgtype и отдаем клиенту
    return c.JSON(http.StatusOK, map[string]interface{}{
        "rating":         statsRow.Rating.Int32,
        "packs_created":  statsRow.PacksCreated.Int32,
        "packs_mastered": statsRow.PacksMastered.Int32,
    })
}


// Вспомогательная функция для преобразования строки в pgtype.UUID
func uuidFromString(s string) pgtype.UUID {
    var u pgtype.UUID
    _ = u.Scan(s) // ошибку игнорируем, при неправильном id придёт 400 на Scan выше
    return u
}
