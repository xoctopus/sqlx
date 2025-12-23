@def T
@def context.Context
@def github.com/xoctopus/sqlx/pkg/builder.SqlCondition
@def github.com/xoctopus/sqlx/pkg/builder.Additions
@def github.com/xoctopus/sqlx/pkg/builder.Cols
@def github.com/xoctopus/sqlx/pkg/builder.Select
@def github.com/xoctopus/sqlx/pkg/builder.And
@def github.com/xoctopus/sqlx/pkg/builder.CC_DriverValue
@def github.com/xoctopus/sqlx/pkg/builder.Eq_DriverValue
@def github.com/xoctopus/sqlx/pkg/helper.Scan
@def github.com/xoctopus/sqlx/pkg/session.For
@def github.com/xoctopus/sqlx/pkg/types.SoftDeletion
--List
func (m *#T#) List(ctx #context.Context#, cond #github.com/xoctopus/sqlx/pkg/builder.SqlCondition#, adds #github.com/xoctopus/sqlx/pkg/builder.Additions#, ignores ...#github.com/xoctopus/sqlx/pkg/builder.Col#) ([]#T#, error) {
	cols := T#T#.TrimmedColumns(ignores...)

	if x, ok := any(m).(#github.com/xoctopus/sqlx/pkg/types.SoftDeletion#); ok {
		deletion, v := x.SoftDeletion()
		cond = #github.com/xoctopus/sqlx/pkg/builder.And#(
			cond,
			#github.com/xoctopus/sqlx/pkg/builder.CC_DriverValue#(T#T#.C(deletion)).AsCond(#github.com/xoctopus/sqlx/pkg/builder.Eq_DriverValue#(v)),
		)
	}
	adds = append(adds, #github.com/xoctopus/sqlx/pkg/builder.Where#(cond))

	rows, err := #github.com/xoctopus/sqlx/pkg/session.For#(ctx, m).Adaptor().Query(
		ctx,
		#github.com/xoctopus/sqlx/pkg/builder.Select#(cols).From(TOrder, adds...),
	)
	if err != nil {
		return nil, err
	}
	res := []#T#{}
	if err = #github.com/xoctopus/sqlx/pkg/helper.Scan#(ctx, rows, &res); err != nil {
		return nil, err
	}
	return res, nil
}
