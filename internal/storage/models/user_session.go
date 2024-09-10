// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// UserSession is an object representing the database table.
type UserSession struct {
	UserID         int64       `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	ChatID         int64       `boil:"chat_id" json:"chat_id" toml:"chat_id" yaml:"chat_id"`
	LastUserAction null.Time   `boil:"last_user_action" json:"last_user_action,omitempty" toml:"last_user_action" yaml:"last_user_action,omitempty"`
	Data           null.String `boil:"data" json:"data,omitempty" toml:"data" yaml:"data,omitempty"`

	R *userSessionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userSessionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserSessionColumns = struct {
	UserID         string
	ChatID         string
	LastUserAction string
	Data           string
}{
	UserID:         "user_id",
	ChatID:         "chat_id",
	LastUserAction: "last_user_action",
	Data:           "data",
}

var UserSessionTableColumns = struct {
	UserID         string
	ChatID         string
	LastUserAction string
	Data           string
}{
	UserID:         "user_session.user_id",
	ChatID:         "user_session.chat_id",
	LastUserAction: "user_session.last_user_action",
	Data:           "user_session.data",
}

// Generated where

var UserSessionWhere = struct {
	UserID         whereHelperint64
	ChatID         whereHelperint64
	LastUserAction whereHelpernull_Time
	Data           whereHelpernull_String
}{
	UserID:         whereHelperint64{field: "\"user_session\".\"user_id\""},
	ChatID:         whereHelperint64{field: "\"user_session\".\"chat_id\""},
	LastUserAction: whereHelpernull_Time{field: "\"user_session\".\"last_user_action\""},
	Data:           whereHelpernull_String{field: "\"user_session\".\"data\""},
}

// UserSessionRels is where relationship names are stored.
var UserSessionRels = struct {
}{}

// userSessionR is where relationships are stored.
type userSessionR struct {
}

// NewStruct creates a new relationship struct
func (*userSessionR) NewStruct() *userSessionR {
	return &userSessionR{}
}

// userSessionL is where Load methods for each relationship are stored.
type userSessionL struct{}

var (
	userSessionAllColumns            = []string{"user_id", "chat_id", "last_user_action", "data"}
	userSessionColumnsWithoutDefault = []string{"chat_id"}
	userSessionColumnsWithDefault    = []string{"user_id", "last_user_action", "data"}
	userSessionPrimaryKeyColumns     = []string{"user_id"}
	userSessionGeneratedColumns      = []string{"user_id"}
)

type (
	// UserSessionSlice is an alias for a slice of pointers to UserSession.
	// This should almost always be used instead of []UserSession.
	UserSessionSlice []*UserSession
	// UserSessionHook is the signature for custom UserSession hook methods
	UserSessionHook func(context.Context, boil.ContextExecutor, *UserSession) error

	userSessionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userSessionType                 = reflect.TypeOf(&UserSession{})
	userSessionMapping              = queries.MakeStructMapping(userSessionType)
	userSessionPrimaryKeyMapping, _ = queries.BindMapping(userSessionType, userSessionMapping, userSessionPrimaryKeyColumns)
	userSessionInsertCacheMut       sync.RWMutex
	userSessionInsertCache          = make(map[string]insertCache)
	userSessionUpdateCacheMut       sync.RWMutex
	userSessionUpdateCache          = make(map[string]updateCache)
	userSessionUpsertCacheMut       sync.RWMutex
	userSessionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userSessionAfterSelectHooks []UserSessionHook

var userSessionBeforeInsertHooks []UserSessionHook
var userSessionAfterInsertHooks []UserSessionHook

var userSessionBeforeUpdateHooks []UserSessionHook
var userSessionAfterUpdateHooks []UserSessionHook

var userSessionBeforeDeleteHooks []UserSessionHook
var userSessionAfterDeleteHooks []UserSessionHook

var userSessionBeforeUpsertHooks []UserSessionHook
var userSessionAfterUpsertHooks []UserSessionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserSession) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserSession) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserSession) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserSession) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserSession) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserSession) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserSession) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserSession) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserSession) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userSessionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserSessionHook registers your hook function for all future operations.
func AddUserSessionHook(hookPoint boil.HookPoint, userSessionHook UserSessionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		userSessionAfterSelectHooks = append(userSessionAfterSelectHooks, userSessionHook)
	case boil.BeforeInsertHook:
		userSessionBeforeInsertHooks = append(userSessionBeforeInsertHooks, userSessionHook)
	case boil.AfterInsertHook:
		userSessionAfterInsertHooks = append(userSessionAfterInsertHooks, userSessionHook)
	case boil.BeforeUpdateHook:
		userSessionBeforeUpdateHooks = append(userSessionBeforeUpdateHooks, userSessionHook)
	case boil.AfterUpdateHook:
		userSessionAfterUpdateHooks = append(userSessionAfterUpdateHooks, userSessionHook)
	case boil.BeforeDeleteHook:
		userSessionBeforeDeleteHooks = append(userSessionBeforeDeleteHooks, userSessionHook)
	case boil.AfterDeleteHook:
		userSessionAfterDeleteHooks = append(userSessionAfterDeleteHooks, userSessionHook)
	case boil.BeforeUpsertHook:
		userSessionBeforeUpsertHooks = append(userSessionBeforeUpsertHooks, userSessionHook)
	case boil.AfterUpsertHook:
		userSessionAfterUpsertHooks = append(userSessionAfterUpsertHooks, userSessionHook)
	}
}

// One returns a single userSession record from the query.
func (q userSessionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserSession, error) {
	o := &UserSession{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_session")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserSession records from the query.
func (q userSessionQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserSessionSlice, error) {
	var o []*UserSession

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserSession slice")
	}

	if len(userSessionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserSession records in the query.
func (q userSessionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_session rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userSessionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_session exists")
	}

	return count > 0, nil
}

// UserSessions retrieves all the records using an executor.
func UserSessions(mods ...qm.QueryMod) userSessionQuery {
	mods = append(mods, qm.From("\"user_session\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"user_session\".*"})
	}

	return userSessionQuery{q}
}

// FindUserSession retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserSession(ctx context.Context, exec boil.ContextExecutor, userID int64, selectCols ...string) (*UserSession, error) {
	userSessionObj := &UserSession{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_session\" where \"user_id\"=?", sel,
	)

	q := queries.Raw(query, userID)

	err := q.Bind(ctx, exec, userSessionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_session")
	}

	if err = userSessionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return userSessionObj, err
	}

	return userSessionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserSession) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_session provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userSessionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userSessionInsertCacheMut.RLock()
	cache, cached := userSessionInsertCache[key]
	userSessionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userSessionAllColumns,
			userSessionColumnsWithDefault,
			userSessionColumnsWithoutDefault,
			nzDefaults,
		)
		wl = strmangle.SetComplement(wl, userSessionGeneratedColumns)

		cache.valueMapping, err = queries.BindMapping(userSessionType, userSessionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userSessionType, userSessionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_session\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_session\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into user_session")
	}

	if !cached {
		userSessionInsertCacheMut.Lock()
		userSessionInsertCache[key] = cache
		userSessionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UserSession.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserSession) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userSessionUpdateCacheMut.RLock()
	cache, cached := userSessionUpdateCache[key]
	userSessionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userSessionAllColumns,
			userSessionPrimaryKeyColumns,
		)
		wl = strmangle.SetComplement(wl, userSessionGeneratedColumns)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update user_session, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_session\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, userSessionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userSessionType, userSessionMapping, append(wl, userSessionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update user_session row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for user_session")
	}

	if !cached {
		userSessionUpdateCacheMut.Lock()
		userSessionUpdateCache[key] = cache
		userSessionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userSessionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for user_session")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for user_session")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserSessionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_session\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userSessionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in userSession slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all userSession")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserSession) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_session provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userSessionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userSessionUpsertCacheMut.RLock()
	cache, cached := userSessionUpsertCache[key]
	userSessionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userSessionAllColumns,
			userSessionColumnsWithDefault,
			userSessionColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userSessionAllColumns,
			userSessionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert user_session, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userSessionPrimaryKeyColumns))
			copy(conflict, userSessionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"user_session\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userSessionType, userSessionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userSessionType, userSessionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert user_session")
	}

	if !cached {
		userSessionUpsertCacheMut.Lock()
		userSessionUpsertCache[key] = cache
		userSessionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UserSession record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserSession) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UserSession provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userSessionPrimaryKeyMapping)
	sql := "DELETE FROM \"user_session\" WHERE \"user_id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from user_session")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for user_session")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userSessionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no userSessionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from user_session")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_session")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserSessionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userSessionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_session\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userSessionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from userSession slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_session")
	}

	if len(userSessionAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserSession) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserSession(ctx, exec, o.UserID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSessionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserSessionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_session\".* FROM \"user_session\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userSessionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserSessionSlice")
	}

	*o = slice

	return nil
}

// UserSessionExists checks if the UserSession row exists.
func UserSessionExists(ctx context.Context, exec boil.ContextExecutor, userID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_session\" where \"user_id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, userID)
	}
	row := exec.QueryRowContext(ctx, sql, userID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_session exists")
	}

	return exists, nil
}
