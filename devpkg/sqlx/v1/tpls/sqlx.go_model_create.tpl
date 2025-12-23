@def T
@def context.Context
@def helper.ColumnsAndValuesForInsertion
@def session.For
@def builder.Insert
@def builder.Comment
--Create
// Create inserts #T# to database
func (m *#T#) Create(ctx #context.Context#) error {
	#CreationMarker#

	cols, values := #helper.ColumnsAndValuesForInsertion#(m)
	_, err := #session.For#(ctx, m).Adaptor().Exec(
		ctx,
		#builder.Insert#().Into(
			T#T#,
			#builder.Comment#(#CreateComment#),
		).Values(cols, values...),
	)
	return err
}

