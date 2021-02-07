module Main exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes as Attr
import Html.Events exposing (onClick)
import Http
import Json.Decode as Decode exposing (Decoder)
import Json.Encode as Encode


main : Program () Model Msg
main =
  Browser.element
    { init = init
    , view = view
    , update = update
    , subscriptions = \_ -> Sub.none
    }

-- MODEL

type Model
  = Failure
  | Loading
  | Success Data

type alias ChartOptions =
    { title: TitleOption
    , tooltip: TooltipOption
    , legend: List LegendOption
    , animationDuration: Int
    , animationEasingUpdate: String
    , series: List SeriesOption
    }

type alias TitleOption =
    { text: String
    , subtext: String
    , top: String
    , left: String
    }

type alias TooltipOption = {}

type alias LegendOption = { data: List String }

type alias SeriesOption =
    { name: String
    , type_: String
    , layout: String
    , roam: Bool
    , label: LabelOption
    , lineStyle: LineStyleOption
    , emphasis: EmphasisOption
    , data: List GraphNodeItemOption
    , links: List GraphEdgeItemOption
    , categories: List CategoryItemOption
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

type alias GraphNodeItemOption =
    { id: String
    , name: String
    , value: Float
    , x: Float
    , y: Float
    , category: Int
    , symbolSize: Float
    , label: Maybe LabelOption
    }

type alias GraphEdgeItemOption =
    { source: String
    , target: String
    }

type alias CategoryItemOption = { name: String }

init : () -> ( Model, Cmd Msg )
init _ = ( Loading, getData )


-- UPDATE

type Msg
    = Fetch
    | GotData (Result Http.Error Data)

update : Msg -> Model -> ( Model, Cmd Msg )
update msg _ =
    case msg of
        Fetch ->
            (Loading, getData)

        GotData result ->
            case result of
                Ok data ->
                  (Success data, Cmd.none)

                Err _ ->
                  (Failure, Cmd.none)


-- VIEW


view : Model -> Html Msg
view model =
    div
        [ Attr.style "display" "flex"
        ]
        [ div
            [ Attr.id "graph"
            , Attr.style "width" "1200px"
            , Attr.style "height" "800px"
            ]
            [ viewGraph model
            ]
        ]

viewGraph : Model -> Html Msg
viewGraph model =
    case model of
        Failure ->
            div []
                [ text "I could not load for some reason. "
                , button [ onClick Fetch ] [ text "Try Again!" ]
                ]
        Loading ->
            text "Loading..."
        Success data ->
            graph (dataToChartOpts data) []

graph : ChartOptions -> List (Html msg) -> Html msg
graph chartOptions =
    node "echart-element"[Attr.property "option" <| encodeChartOptions chartOptions ]


dataToChartOpts : Data -> ChartOptions
dataToChartOpts data =
    { title =
        { text = "Les Miserables"
        , subtext = "Default layout"
        , top = "bottom"
        , left = "right"
        }
    , tooltip = {}
    , legend = [{ data = List.map (\d -> d.name ) data.categories }]
    , animationDuration = 1500
    , animationEasingUpdate = "quinticInOut"
    , series = [
        { name = "Les Miserables"
        , type_ = "graph"
        , layout = "none"
        , roam = True
        , label =
            { position = Just "right"
            , formatter = Just "{b}"
            , show = Nothing
            }
        , lineStyle =
            { color = Just "source"
            , curveness = Just 0.3
            , width = Nothing
            }
        , emphasis =
            { focus = "adjacency"
            , lineStyle =
                { width = Just 10
                , color = Nothing
                , curveness = Nothing
                }
            }
        , data = List.map setLabel data.nodes
        , links = data.links
        , categories = data.categories
        }]
    } |> Debug.log "ToChartOps: "

setLabel : GraphNodeItemOption -> GraphNodeItemOption
setLabel graphNodeItemOption =
    { graphNodeItemOption | label =
        Just { show = Just(graphNodeItemOption.symbolSize > 30)
             , formatter = Nothing
             , position = Nothing
             }
    }

-- HTTP


getData : Cmd Msg
getData =
  Http.get
    { url = "https://cors-anywhere.herokuapp.com/https://echarts.apache.org/next/examples/data/asset/data/les-miserables.json"
    , expect = Http.expectJson GotData dataDecoder
    }

type alias Data =
    { nodes: List GraphNodeItemOption
    , links: List GraphEdgeItemOption
    , categories: List CategoryItemOption
    }

dataDecoder : Decoder Data
dataDecoder =
    Decode.map3 Data
        (Decode.field "nodes" (graphNodeItemOptionDecoder))
        (Decode.field "links" (graphEdgeItemOptionDecoder))
        (Decode.field "categories" (categoryItemOptionDecoder))

graphNodeItemOptionDecoder : Decoder (List GraphNodeItemOption)
graphNodeItemOptionDecoder =
    Decode.map8 GraphNodeItemOption
        (Decode.field "id" Decode.string)
        (Decode.field "name" Decode.string)
        (Decode.field "value" Decode.float)
        (Decode.field "x" Decode.float)
        (Decode.field "y" Decode.float)
        (Decode.field "category" Decode.int)
        (Decode.field "symbolSize" Decode.float)
        (Decode.maybe (Decode.field "label" labelOptionDecoder))
        |> Decode.list

labelOptionDecoder : Decoder LabelOption
labelOptionDecoder =
    Decode.map3 LabelOption
        (Decode.maybe (Decode.field "position" Decode.string))
        (Decode.maybe (Decode.field "formatter" Decode.string))
        (Decode.maybe (Decode.field "show" Decode.bool))


graphEdgeItemOptionDecoder : Decoder (List GraphEdgeItemOption)
graphEdgeItemOptionDecoder =
    Decode.map2 GraphEdgeItemOption
        (Decode.field "source" Decode.string)
        (Decode.field "target" Decode.string)
        |> Decode.list


categoryItemOptionDecoder : Decoder (List CategoryItemOption)
categoryItemOptionDecoder =
    Decode.map CategoryItemOption
        (Decode.field "name" Decode.string)
        |> Decode.list


encodeChartOptions : ChartOptions -> Encode.Value
encodeChartOptions chartOptions =
    Encode.object
        [ ( "title", titleOptionEncoder chartOptions.title )
        , ( "tooltip", tooltipOptionEncoder )
        , ( "legend", Encode.list legendOptionEncoder chartOptions.legend )
        , ( "animationDuration", Encode.int chartOptions.animationDuration )
        , ( "animationEasingUpdate", Encode.string chartOptions.animationEasingUpdate )
        , ( "series", Encode.list seriesOptionEncoder chartOptions.series )
        ]

titleOptionEncoder : TitleOption -> Encode.Value
titleOptionEncoder titleOption =
    Encode.object
        [ ( "text", Encode.string titleOption.text )
        , ( "subtext", Encode.string titleOption.subtext )
        , ( "top", Encode.string titleOption.top )
        , ( "left", Encode.string titleOption.left )
        ]

tooltipOptionEncoder : Encode.Value
tooltipOptionEncoder = Encode.object []

legendOptionEncoder : LegendOption -> Encode.Value
legendOptionEncoder legendOption = Encode.object [ ( "data", Encode.list Encode.string legendOption.data ) ]

seriesOptionEncoder : SeriesOption -> Encode.Value
seriesOptionEncoder seriesOption =
    Encode.object
        [ ( "name", Encode.string seriesOption.name )
        , ( "type", Encode.string seriesOption.type_ )
        , ( "layout", Encode.string seriesOption.layout )
        , ( "roam", Encode.bool seriesOption.roam )
        , ( "label", labelOptionEncoder seriesOption.label )
        , ( "lineStyle", lineStyleOptionEncoder seriesOption.lineStyle )
        , ( "emphasis", emphasisOptionEncoder seriesOption.emphasis )
        , ( "data", Encode.list graphNodeItemOptionEncoder seriesOption.data )
        , ( "links", Encode.list graphEdgeItemOptionEncoder seriesOption.links )
        , ( "categories", Encode.list categoryItemOptionEncoder seriesOption.categories )
        ]

labelOptionEncoder : LabelOption -> Encode.Value
labelOptionEncoder labelOption =
    (Encode.object << List.filterMap identity)
        [ labelOption.position
            |> Maybe.andThen (\position -> Just ( "position", Encode.string position ))
        , labelOption.formatter
            |> Maybe.andThen (\formatter -> Just ( "formatter", Encode.string formatter ))
        , labelOption.show
            |> Maybe.andThen (\show -> Just ( "show", Encode.bool show ))
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

graphNodeItemOptionEncoder : GraphNodeItemOption -> Encode.Value
graphNodeItemOptionEncoder graphNodeItemOption =
    (Encode.object << List.filterMap identity)
        [ Just( "id", Encode.string graphNodeItemOption.id )
        , Just( "name", Encode.string graphNodeItemOption.name )
        , Just( "value", Encode.float graphNodeItemOption.value )
        , Just( "x", Encode.float graphNodeItemOption.x )
        , Just( "y", Encode.float graphNodeItemOption.y )
        , Just( "category", Encode.int graphNodeItemOption.category )
        , Just( "symbolSize", Encode.float graphNodeItemOption.symbolSize )
        , graphNodeItemOption.label |> Maybe.andThen (\label -> Just ("label", labelOptionEncoder label))
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
