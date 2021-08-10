module EntityHelpers exposing (..)

import Entity.Object
import Entity.Object.Entity as Entity
import Entity.Query
import Entity.Scalar exposing (Entity(..))
import Entity.ScalarCodecs exposing (Id)
import Graphql.Operation exposing (RootQuery)
import Graphql.OptionalArgument as OptionalArgument
import Graphql.SelectionSet as SelectionSet exposing (SelectionSet, with)


type alias EntityListResponse =
    List EntityResponse


type alias EntityResponse =
    Entity


type alias Entity =
    { id : Id
    , name : String
    , context : String
    , connections : List String
    }


entitySelection : SelectionSet Entity Entity.Object.Entity
entitySelection =
    SelectionSet.succeed Entity
        |> with Entity.id
        |> with Entity.name
        |> with Entity.context
        |> with (Entity.connections (\_ -> { id = OptionalArgument.Absent, context = OptionalArgument.Absent }) Entity.name)


listEntitiesQuery : SelectionSet EntityListResponse RootQuery
listEntitiesQuery =
    Entity.Query.list entitySelection


getEntityQuery : String -> SelectionSet (Maybe EntityResponse) RootQuery
getEntityQuery id =
    Entity.Query.entity { id = Id id } entitySelection


toString : Id -> String
toString id =
    case id of
        Id string ->
            string
