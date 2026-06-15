@def T
@def context.Context
@def helper.CVsForInsertion
@def session.MustExecutorFor
@def builder.Insert
@def builder.Comment
--Create
// Create inserts #T# to database
func (m *#T#) Create(ctx #context.Context#) error {
	#CreationMarker#

	cols, values := #helper.CVsForInsertion#(m)
	_, err := #session.MustExecutorFor#(ctx, T#T#).Exec(
		ctx,
		#builder.Insert#().Into(
			T#T#,
			#builder.Comment#(#CreateComment#),
		).Values(cols, values...),
	)
	return err
}

