module AdminHelpers exposing (..)

import Admin.Query
import Graphql.Operation exposing (RootQuery)
import Graphql.SelectionSet exposing (SelectionSet)
type alias AdminResponse =
    Maybe String

rebuildQuery : SelectionSet AdminResponse RootQuery
rebuildQuery = Admin.Query.rebuild
