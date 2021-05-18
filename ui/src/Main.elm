module Main exposing (main)

import Admin.Query as AdminQuery
import Browser
import Css
import Debug
import Entity.Object
import Entity.Object.Entity as Entity
import Entity.Query as EntityQuery
import Entity.Scalar exposing (Id(..))
import Graphql.Http exposing (Error, HttpError(..))
import Graphql.Http.GraphqlError exposing (GraphqlError, PossiblyParsedData(..))
import Graphql.Operation exposing (RootQuery)
import Graphql.OptionalArgument as OptionalArgument
import Graphql.SelectionSet as SelectionSet exposing (SelectionSet, with)
import Html.Styled.Events exposing (onClick)
import Http
import Json.Encode as Encode
import List
import List.Extra
import Material.Icons as Icons
import Material.Icons.Types exposing (Coloring(..), Icon)
import RemoteData exposing (RemoteData(..))
import Entity.Object.Entity as Entity
import Svg.Styled as Svg
import Svg.Styled.Attributes exposing (css, d, r, viewBox)
import Tailwind.Utilities as Tw
import Css.Global
import Html.Styled as Html
import Html.Styled.Attributes as Attr

-- INBOUND TYPES

type alias Data =
    { nodes: List GraphNodeItemOption
    , links: List GraphEdgeItemOption
    , categories: List CategoryItemOption
    }

type alias GraphEdgeItemOption =
    { source: String
    , target: String
    }

type alias GraphNodeItemOption =
    { name: Maybe String
    , value: Float
    , x: Maybe Float
    , y: Maybe Float
    , category: Maybe Int
    , symbolSize: Maybe Float
    , label: Maybe LabelOption
    }

type alias CategoryItemOption =
    { name: String
    }

type alias Entity =
    { id : Id
    , name : String
    , context : String
    , connections: List String
    }

type alias AdminResponse =
    Maybe String

type alias EntityResponse =
    Entity

type alias EntityListResponse =
    List EntityResponse


-- TARGET TYPES

type alias ChartOptions =
    { title: Maybe TitleOption
    , tooltip: Maybe TooltipOption
    , legend: List LegendOption
    , series: List SeriesOption
    }

type alias TitleOption =
    { text: String
    , subtext: String
    , top: String
    , left: String
    }

type alias TooltipOption =
    { show: Bool
    }

type alias LegendOption =
    { data: List String
    }

type alias SeriesOption =
    { animation: Maybe Bool
    , categories: List CategoryItemOption
    , data: List GraphNodeItemOption
    , draggable: Maybe Bool
    , emphasis: Maybe EmphasisOption
    , force: Maybe ForceOption
    , label: Maybe LabelOption
    , layout: String
    , lineStyle: Maybe LineStyleOption
    , links: List GraphEdgeItemOption
    , name: Maybe String
    , roam: Maybe Bool
    , type_: String
    }

type alias LabelOption =
    { position: Maybe String
    , formatter: Maybe String
    , show: Maybe Bool
    }

type alias LineStyleOption =
    { color: Maybe String
    , curveness: Maybe Float
    , width: Maybe Int
    }

type alias EmphasisOption =
    { focus: String
    , lineStyle: LineStyleOption
    }

type alias ForceOption =
    { edgeLength: Maybe Int
    , friction: Maybe Float
    , gravity: Maybe Float
    , layoutAnimation: Maybe Bool
    , repulsion: Maybe Int
    }

entitySelection : SelectionSet Entity Entity.Object.Entity
entitySelection =
    SelectionSet.succeed Entity
        |> with Entity.id
        |> with Entity.name
        |> with Entity.context
        |> with (Entity.connections (\_ -> {id = OptionalArgument.Absent, context = OptionalArgument.Absent} ) Entity.name)


-- ENCODERS

encodeChartOptions : ChartOptions -> Encode.Value
encodeChartOptions chartOptions =
    (Encode.object << List.filterMap identity)
        [ chartOptions.title |> Maybe.andThen (\title -> Just ( "title", titleOptionEncoder title ))
        , chartOptions.tooltip |> Maybe.andThen (\t -> Just ( "tooltip", tooltipOptionEncoder t ))
        , Just( "legend", Encode.list legendOptionEncoder chartOptions.legend)
        , Just( "series", Encode.list seriesOptionEncoder chartOptions.series )
        ]

titleOptionEncoder : TitleOption -> Encode.Value
titleOptionEncoder titleOption =
    Encode.object
        [ ( "text", Encode.string titleOption.text )
        , ( "subtext", Encode.string titleOption.subtext )
        , ( "top", Encode.string titleOption.top )
        , ( "left", Encode.string titleOption.left )
        ]

tooltipOptionEncoder : TooltipOption -> Encode.Value
tooltipOptionEncoder tooltipOption =
    Encode.object
        [ ("show", Encode.bool tooltipOption.show )
        ]

legendOptionEncoder : LegendOption -> Encode.Value
legendOptionEncoder legendOption = Encode.object [ ( "data", Encode.list Encode.string legendOption.data ) ]

seriesOptionEncoder : SeriesOption -> Encode.Value
seriesOptionEncoder seriesOption =
    (Encode.object << List.filterMap identity)
        [ seriesOption.animation |> Maybe.andThen (\a -> Just ("animation", Encode.bool a))
        , Just ( "categories", Encode.list categoryItemOptionEncoder seriesOption.categories )
        , Just ( "data", Encode.list graphNodeItemOptionEncoder seriesOption.data )
        , seriesOption.draggable |> Maybe.andThen (\d -> Just ( "draggable", Encode.bool d))
        , seriesOption.emphasis |> Maybe.andThen (\e -> Just ( "emphasis", emphasisOptionEncoder e ))
        , seriesOption.force |> Maybe.andThen (\f -> Just ("force", forceOptionEncoder f ))
        , seriesOption.label |> Maybe.andThen (\l -> Just ( "label", labelOptionEncoder l ))
        , Just ( "layout", Encode.string seriesOption.layout )
        , seriesOption.lineStyle |> Maybe.andThen (\ls -> Just ( "lineStyle", lineStyleOptionEncoder ls ))
        , Just ( "links", Encode.list graphEdgeItemOptionEncoder seriesOption.links )
        , seriesOption.name |> Maybe.andThen (\n -> Just ( "name", Encode.string n ))
        , seriesOption.roam |> Maybe.andThen (\r -> Just ("roam", Encode.bool r ))
        , Just ( "type", Encode.string seriesOption.type_ )
        ]

labelOptionEncoder : LabelOption -> Encode.Value
labelOptionEncoder labelOption =
    (Encode.object << List.filterMap identity)
        [ labelOption.position |> Maybe.andThen (\p -> Just ( "position", Encode.string p ))
        , labelOption.formatter |> Maybe.andThen (\f -> Just ( "formatter", Encode.string f ))
        , labelOption.show |> Maybe.andThen (\s -> Just ( "show", Encode.bool s ))
        ]

lineStyleOptionEncoder : LineStyleOption -> Encode.Value
lineStyleOptionEncoder lineStyleOption =
    (Encode.object << List.filterMap identity)
        [ lineStyleOption.color
            |> Maybe.andThen (\color -> Just ( "color", Encode.string color ))
        , lineStyleOption.curveness
            |> Maybe.andThen (\curveness -> Just ( "curveness", Encode.float curveness ))
        , lineStyleOption.width
            |> Maybe.andThen (\width -> Just ( "width", Encode.int width ))
        ]

emphasisOptionEncoder : EmphasisOption -> Encode.Value
emphasisOptionEncoder emphasisOption =
    Encode.object
        [ ( "focus", Encode.string emphasisOption.focus )
        , ( "lineStyle", lineStyleOptionEncoder emphasisOption.lineStyle)
        ]

forceOptionEncoder : ForceOption -> Encode.Value
forceOptionEncoder forceOption =
    (Encode.object << List.filterMap identity)
        [ forceOption.edgeLength |> Maybe.andThen (\el -> Just ("edgeLength", Encode.int el))
        , forceOption.friction |> Maybe.andThen (\f -> Just ("friction", Encode.float f))
        , forceOption.gravity |> Maybe.andThen (\g -> Just ("gravity", Encode.float g))
        , forceOption.layoutAnimation |> Maybe.andThen (\la -> Just ("layoutAnimation", Encode.bool la))
        , forceOption.repulsion |> Maybe.andThen (\r -> Just ("repulsion", Encode.int r))
        ]

graphNodeItemOptionEncoder : GraphNodeItemOption -> Encode.Value
graphNodeItemOptionEncoder gnio =
    (Encode.object << List.filterMap identity)
        [ gnio.name |> Maybe.andThen (\n -> Just( "name", Encode.string n ))
        , Just( "value", Encode.float gnio.value )
        , gnio.x |> Maybe.andThen (\x -> Just ( "x", Encode.float x ))
        , gnio.y |> Maybe.andThen (\y -> Just ( "y", Encode.float y ))
        , gnio.category |> Maybe.andThen (\c -> Just( "category", Encode.int c ))
        , gnio.symbolSize |> Maybe.andThen (\s -> Just ( "symbolSize", Encode.float s ))
        , gnio.label |> Maybe.andThen (\l -> Just ("label", labelOptionEncoder l))
        ]

graphEdgeItemOptionEncoder : GraphEdgeItemOption -> Encode.Value
graphEdgeItemOptionEncoder graphEdgeItemOption =
    Encode.object
        [ ( "source", Encode.string graphEdgeItemOption.source )
        , ( "target", Encode.string graphEdgeItemOption.target )
        ]

categoryItemOptionEncoder : CategoryItemOption -> Encode.Value
categoryItemOptionEncoder categoryItemOption =
    Encode.object
        [ ( "name", Encode.string categoryItemOption.name ) ]

-- FUNCTIONS

entityListToChartOpts : EntityListResponse -> ChartOptions
entityListToChartOpts entityList =
    let
        categories = List.sort (List.Extra.unique (List.map (\e -> e.context) entityList))
    in
        { title = Nothing
        , legend = [{ data = categories }]
        , tooltip =
            Just { show = False
            }
        , series = [
            { animation = Just True
            , categories = List.map (\c -> { name = c }) categories
            , data = List.map (entityToGraphNodeItemOption categories) entityList
            , draggable = Just True
            , emphasis = Nothing
            , force =
                Just
                    { edgeLength = Just 100
                    , friction = Just 0.2
                    , gravity = Just 0.2
                    , layoutAnimation = Just True
                    , repulsion = Just 50
                    }
            , label =
                Just
                    { position = Just "right"
                    , formatter = Just "{b}"
                    , show = Nothing
                    }
            , layout = "force"
            , lineStyle = Nothing
            , links = List.concatMap entityToGraphEdgeItemOption entityList
            , name = Nothing
            , roam = Just True
            , type_ = "graph"
            }]
        }


entityToGraphNodeItemOption : List String -> Entity -> GraphNodeItemOption
entityToGraphNodeItemOption l e =
    { name = Just e.name
    , category = List.Extra.elemIndex e.context l
    , value = 1
    , x = Nothing
    , y = Nothing
    , symbolSize = Just 20
    , label = Nothing
    }

entityToGraphEdgeItemOption: Entity -> List GraphEdgeItemOption
entityToGraphEdgeItemOption e =
    List.map (\c -> { source = e.name, target = c}) e.connections


-- GRAPHQL


rebuildQuery : SelectionSet AdminResponse RootQuery
rebuildQuery = AdminQuery.rebuild

listEntitiesQuery : SelectionSet EntityListResponse RootQuery
listEntitiesQuery = EntityQuery.list ( entitySelection )

getEntityQuery : String -> SelectionSet (Maybe EntityResponse) RootQuery
getEntityQuery id =
    EntityQuery.entity { id = (Id id) } entitySelection


makeListEntitiesQuery : Cmd Msg
makeListEntitiesQuery =
    listEntitiesQuery
        |> Graphql.Http.queryRequest "http://127.0.0.1:5000/entity"
        |> Graphql.Http.send
            (Graphql.Http.discardParsedErrorData
                >> RemoteData.fromResult
                >> GotEntityListResponse
            )

sendRebuildQuery : Cmd Msg
sendRebuildQuery =
    rebuildQuery
        |> Graphql.Http.queryRequest "http://127.0.0.1:5000/admin"
        |> Graphql.Http.send
            (Graphql.Http.discardParsedErrorData
                >> RemoteData.fromResult
                >> GotAdminResponse
            )

--drawGraph : Cmd Msg
--drawGraph =



-- ELM ARCHITECTURE

main : Program () Model Msg
main =
  Browser.element
    { init = init
    , view = view >> Html.toUnstyled
    , update = update
    , subscriptions = \_ -> Sub.none
    }

 -- INIT

init : () -> ( Model, Cmd Msg )
init _ =
    ( { entities = RemoteData.Loading
      , indexing = RemoteData.NotAsked
      }
    , makeListEntitiesQuery
    )

-- MODEL

type alias EntityListModel = RemoteData (Graphql.Http.Error ()) EntityListResponse
type alias AdminModel = RemoteData (Graphql.Http.Error ()) AdminResponse

type alias Model =
    { entities: EntityListModel
    , indexing: AdminModel
    }

-- UPDATE


type Msg
    = GotEntityListResponse EntityListModel
    | FetchEntityList
    | GotAdminResponse AdminModel
    | FetchAdminResponse
    --| Draw ChartOptions


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotEntityListResponse entityListModel ->
            ( { model | entities = entityListModel }
            , Cmd.none
            )
        FetchEntityList ->
            ( { model | entities = RemoteData.Loading }
            , makeListEntitiesQuery
            )
        GotAdminResponse adminModel ->
            ( { model | indexing = adminModel }
            , Cmd.none
            )
        FetchAdminResponse ->
            ( { model | indexing = RemoteData.Loading }
            , sendRebuildQuery
            )
        --Draw chartOptions ->
        --    ( model, drawGraph chartOptions)


-- VIEW


view : Model -> Html.Html Msg
view model =
    Html.div
        [ Attr.css
            [ Tw.bg_gray_50 ]
        , Attr.style "display" "flex"
        ]
        [ Css.Global.global Tw.globalStyles
        , Html.div
            [ Attr.id "graph"
            , Attr.style "width" "1200px"
            , Attr.style "height" "800px"
            ]
            [ viewEntityListModelResult model.entities
            ]
        , Html.div
            [ Attr.id "rebuildBtn" ]
            [ viewAdminModelResult model.indexing ]
        ]

viewEntityListModelResult : EntityListModel -> Html.Html Msg
viewEntityListModelResult model =
    case model of
        NotAsked ->
            Html.text "I didn't ask"

        Loading ->
            Html.text "Loading..."

        Failure e ->
            Html.div [] (buildFailureMsg e)

        Success entityList ->
            graph (entityListToChartOpts entityList) []

buildFailureMsg: Error parsedData -> List (Html.Html Msg)
buildFailureMsg parsedData =
    case parsedData of
        Graphql.Http.GraphqlError _ graphqlErrors ->
            List.map (\err -> buildErrorMsg "Graphql Error" err.message) graphqlErrors

        Graphql.Http.HttpError httpError ->
            [ buildErrorMsg "Http Error" (buildHttpErrorMessage httpError) ]

buildErrorMsg: String -> String -> Html.Html Msg
buildErrorMsg eType eMsg =
    Html.div
        [ Attr.css
            [ Tw.flex
            , Tw.bg_red_200
            , Tw.p_4
            ]
        ]
        [ Html.div
            [ Attr.css [ Tw.mr_4 ] ]
            [ Html.div
                [ Attr.css
                    [ Tw.h_10
                    , Tw.w_10
                    , Tw.text_white
                    , Tw.bg_red_600
                    , Tw.rounded_full
                    , Tw.flex
                    , Tw.justify_center
                    , Tw.items_center
                    ]
                ]
                [ Html.fromUnstyled (Icons.warning 24 Inherit) ]
            ]
        , Html.div
            [ Attr.css
                [ Tw.flex
                , Tw.justify_between
                , Tw.w_full
                ]
            ]
            [ Html.div
                [ Attr.css [ Tw.text_red_600 ] ]
                [ Html.p
                    [ Attr.css
                        [ Tw.mb_2
                        , Tw.font_bold
                        ]
                    ]
                    [ Html.text eType ]
                , Html.p [ Attr.css [ Tw.text_xs ] ] [ Html.text eMsg ]
                ]
            , Html.div
                [ Attr.css
                    [ Tw.text_sm
                    , Tw.text_gray_500
                    ]
                ]
                [ Html.button [] [ Html.text "x" ] ]
            ]
        ]


buildHttpErrorMessage : HttpError -> String
buildHttpErrorMessage httpError =
    case httpError of
        BadUrl message ->
            message

        Timeout ->
            "Server is taking too long to respond. Please try again later."

        NetworkError ->
            "Unable to reach server."

        BadStatus metadata body ->
            "Request failed with status code: " ++ String.fromInt metadata.statusCode ++ ". Error: " ++ body

        BadPayload error ->
            "Bad payload received: " ++ Debug.toString error

graph : ChartOptions -> List (Html.Html msg) -> Html.Html msg
graph chartOptions =
    Html.node "echart-element" [ Attr.property "option" <| encodeChartOptions chartOptions ]

viewAdminModelResult : AdminModel -> Html.Html Msg
viewAdminModelResult model =
    case model of
        RemoteData.NotAsked ->
            rebuildBtn
                Tw.bg_blue_500
                (rebuildBtnSvg rebuildBtnLogoSpec Icons.build "Rebuild")


        RemoteData.Loading ->
            rebuildBtn
                Tw.bg_yellow_500
                (rebuildBtnSvg (Tw.animate_spin :: rebuildBtnLogoSpec) Icons.refresh "Building")

        RemoteData.Failure e ->
            rebuildBtn
                Tw.bg_red_600
                (rebuildBtnSvg rebuildBtnLogoSpec Icons.error_outline ("Error: " ++ Debug.toString e) )

        RemoteData.Success _ ->
            rebuildBtn
                Tw.bg_green_400
                (rebuildBtnSvg rebuildBtnLogoSpec Icons.check_circle_outline "Rebuilt")

rebuildBtnLogoSpec : List (Css.Style)
rebuildBtnLogoSpec =
    [ Tw.h_5
    , Tw.w_5
    , Tw.mr_3
    ]

rebuildBtnSvg : List (Css.Style) -> Icon msg -> String -> List (Html.Html msg)
rebuildBtnSvg cssStyle icon displayText =
    [ Svg.svg
        [ css cssStyle
        , viewBox "0 0 24 24"
        ]
        [ Html.fromUnstyled (icon 24 Inherit) ]
    , Html.text displayText
    ]

rebuildBtn : Css.Style -> List (Html.Html Msg) -> Html.Html Msg
rebuildBtn colour inner =
    Html.button
        [ Attr.css
            [ Tw.flex
              , Tw.items_center
              , Tw.shadow
              , Tw.px_4
              , Tw.py_2
              , Tw.text_white
              , Tw.rounded_md
              , colour
              ]
            , onClick FetchAdminResponse
        ]
        inner
