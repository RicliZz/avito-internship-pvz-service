package authRepo

//func TestRegister_Success(t *testing.T) {
//	logger.Logger = zap.NewNop().Sugar()
//	mock, err := pgxmock.NewPool()
//	assert.NoError(t, err)
//	defer mock.Close()
//
//	repo := &AuthRepository{db: mock}
//
//	email := "test@example.com"
//	password := "hashedpassword"
//	role := "user"
//
//	userID := uuid.New()
//
//	rows := pgxmock.NewRows([]string{"ID", "email", "role"}).
//		AddRow(userID, email, role)
//
//	mock.ExpectQuery(`INSERT INTO users \(email, password, role\) VALUES \(\$1, \$2, \$3\) RETURNING "ID", email, role`).
//		WithArgs(email, password, role).
//		WillReturnRows(rows)
//
//	user, err := repo.Register(models.RegisterParams{
//		Email:    email,
//		Password: password,
//		Role:     role,
//	})
//
//	assert.NoError(t, err)
//	assert.Equal(t, userID, user.ID)
//	assert.Equal(t, email, user.Email)
//	assert.Equal(t, role, user.Role)
//
//	assert.NoError(t, mock.ExpectationsWereMet())
//}

//func TestRegister_Fail_DBError(t *testing.T) {
//	logger.Logger = zap.NewNop().Sugar()
//	mock, err := pgxmock.NewPool()
//	assert.NoError(t, err)
//	defer mock.Close()
//
//	repo := &AuthRepository{db: mock}
//
//	email := "test@example.com"
//	password := "hashedpassword"
//	role := "user"
//
//	mock.ExpectQuery(`INSERT INTO users \(email, password, role\) VALUES \(\$1, \$2, \$3\) RETURNING "ID", email, role`).
//		WithArgs(email, password, role).
//		WillReturnError(fmt.Errorf("database error"))
//
//	user, err := repo.Register(models.RegisterParams{
//		Email:    email,
//		Password: password,
//		Role:     role,
//	})
//
//	assert.Error(t, err)
//	assert.Nil(t, user)
//	assert.Equal(t, "database error", err.Error())
//
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
