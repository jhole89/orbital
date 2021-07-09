module Main exposing (main)

import AdminHelpers exposing (AdminResponse, rebuildQuery)
import Array
import Browser
import Css
import ECharts exposing (ChartOptions, GraphEdgeItemOption, GraphNodeItemOption, encodeChartOptions)
import Graphql.Http exposing (Error, HttpError(..))
import Graphql.Http.GraphqlError exposing (GraphqlError, PossiblyParsedData(..))
import EntityHelpers exposing (Entity, EntityListResponse, listEntitiesQuery, toString)
import Html.Styled.Events as Events
import Json.Decode as Decode exposing (errorToString)
import Material.Icons as Icons
import Material.Icons.Types exposing (Coloring(..), Icon)
import RemoteData exposing (RemoteData)
import Svg.Styled as Svg
import Svg.Styled.Attributes as SvgAttr
import Tailwind.Utilities as Tw
import Tailwind.Breakpoints as Bp
import Css.Global exposing (global)
import Html.Styled as Html
import Html.Styled.Attributes as HtmlAttr
import Utilities exposing (entityListToChartOpts)


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
    ( initModelState
    , makeListEntitiesQuery
    )

initModelState : Model
initModelState =
    { entities = RemoteData.Loading
    , indexing = RemoteData.NotAsked
    , chartConfig = NoOpt
    , selectedEntity = NoEnt
    , idReferences = []
    , displaySidebar = True
    , displayWarningModal = False
    }

-- MODEL

type alias EntityListModel = RemoteData (Graphql.Http.Error ()) EntityListResponse
type alias AdminModel = RemoteData (Graphql.Http.Error ()) AdminResponse

type alias IdReferenceIndex = List Entity
type alias SidebarDisplayModel = Bool
type alias WarningDisplayModel = Bool

type ChartOptionsModel =
    ChartConfig ChartOptions
    | NoOpt
type SelectedEntityModel =
    Selected Entity
    | NoEnt

type alias Model =
    { entities: EntityListModel
    , indexing: AdminModel
    , chartConfig: ChartOptionsModel
    , idReferences: IdReferenceIndex
    , selectedEntity: SelectedEntityModel
    , displaySidebar: SidebarDisplayModel
    , displayWarningModal: WarningDisplayModel
    }

-- UPDATE


type Msg =
    GotEntityListResponse EntityListModel
    | FetchAdminResponse
    | GotAdminResponse AdminModel
    | GotId Int
    | SelectEntity SelectedEntityModel
    | ShowSidebar Bool
    | ShowWarningModal Bool


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotEntityListResponse entityModel ->
            case entityModel of
                RemoteData.Success entityListResponse ->
                    ( { model
                        | entities = entityModel
                        , chartConfig = ChartConfig (entityListToChartOpts entityListResponse)
                        , idReferences = entityListResponse
                      }
                      , Cmd.none
                    )
                _ ->
                    ( { model | entities = entityModel }
                    , Cmd.none
                    )
        FetchAdminResponse ->
            ( { model | indexing = RemoteData.Loading }
            , sendRebuildQuery
            )
        GotAdminResponse adminModel ->
            case adminModel of
                RemoteData.Success _ ->
                    ( { model
                        | indexing = adminModel
                        , entities = RemoteData.NotAsked
                      }
                    , makeListEntitiesQuery
                    )
                _ ->
                    ( { model | indexing = adminModel }
                    , Cmd.none
                    )
        GotId value ->
            case model.idReferences |> Array.fromList |> Array.get value of
                Just entity ->
                    ( { model | selectedEntity = Selected entity }
                    , Cmd.none
                    )
                Nothing ->
                    ( { model | selectedEntity = NoEnt }
                    , Cmd.none
                    )
        SelectEntity selectedEntityModel ->
            ( { model | selectedEntity = selectedEntityModel }
            , Cmd.none
            )
        ShowSidebar status ->
            ( { model | displaySidebar = status }
            , Cmd.none
            )
        ShowWarningModal status ->
            ( { model | displayWarningModal = status }
            , Cmd.none
            )


-- VIEW


view : Model -> Html.Html Msg
view model =
    homePage
        [ global Tw.globalStyles
        , mainPage model.displayWarningModal
            [ header model.indexing
            , mainSection
                [ sideBar
                    model.displaySidebar
                    [ sideBarHeaderPanel
                        [ sideBarHeader ]
                    , entityDetails model.selectedEntity
                    , upgradePanel
                    ]
                , sideBarShowBtn model.displaySidebar
                , canvasSection
                    <| viewEntityListModelResult
                    <| model
                ]
            ]
        ]

homePage : List (Html.Html Msg) -> Html.Html Msg
homePage contents =
    Html.div
        [ HtmlAttr.css
            [ Tw.flex
            , Tw.flex_col
            , Bp.md [ Tw.flex_row ]
            ]
        ] contents

mainPage : WarningDisplayModel -> List (Html.Html Msg) -> Html.Html Msg
mainPage model contents =
    case model of
        True ->
            Html.div [][]
        False ->
            Html.div
                [ HtmlAttr.css
                    [ Tw.w_full
                    , Tw.flex
                    , Tw.flex_col
                    , Tw.h_screen
                    , Tw.overflow_y_hidden
                    , Tw.bg_gray_300
                    ]
                ] contents


-- SIDEBAR

sideBar: SidebarDisplayModel -> List (Html.Html Msg) -> Html.Html Msg
sideBar model contents =
    case model of
        True ->
            Html.aside
                [ HtmlAttr.css
                    [ Tw.relative
                    , Tw.bg_white
                    , Tw.rounded_md
                    , Tw.h_full
                    , Tw.w_96
                    , Tw.flex
                    , Tw.hidden
                    , Bp.sm
                        [ Tw.block
                        , Tw.shadow_xl
                        ]
                    ]
                ] contents
        False ->
            Html.aside [][]

sideBarHeaderPanel : List (Html.Html Msg) -> Html.Html Msg
sideBarHeaderPanel contents =
    Html.div
        [ HtmlAttr.css
            [ Tw.flex
            , Tw.bg_blue_500
            , Tw.justify_between
            , Bp.sm [ Tw.px_6 ]
            ]
        ] contents

sideBarHeader : Html.Html Msg
sideBarHeader =
    Html.h2
        [ HtmlAttr.css
            [ Tw.p_6
            , Tw.text_xl
            , Tw.font_bold
            , Tw.text_white
            , Tw.content_center
            , Css.hover [ Tw.text_gray_300 ]
            ]
        ]
        [ Html.text "Telescope" ]

sideBarShowBtn : SidebarDisplayModel -> Html.Html Msg
sideBarShowBtn model =
    case model of
        True ->
            sideBarShowBtnElement (ShowSidebar False) "<"
        False ->
            sideBarShowBtnElement (ShowSidebar True) ">"

sideBarShowBtnElement : Msg -> String -> Html.Html Msg
sideBarShowBtnElement onClickEventMsg displayText =
    Html.button
        [ sideBarShowBtnStyle
        , Events.onClick onClickEventMsg
        ]
        [ Html.text displayText ]

sideBarShowBtnStyle : Html.Attribute Msg
sideBarShowBtnStyle =
    HtmlAttr.css
        [ Tw.bg_blue_700
        , Tw.z_50
        , Tw.mt_1
        , Tw.py_3
        , Tw.px_2
        , Tw.h_16
        , Tw.text_xl
        , Tw.font_bold
        , Tw.text_white
        , Tw.rounded_r_lg
        , Tw.flex
        , Css.hover [ Tw.bg_blue_900 ]
        , Css.focus
            [ Tw.outline_none
            , Tw.ring
            , Tw.ring_blue_500
            ]
        ]

rebuildBtn : AdminModel -> Html.Html Msg
rebuildBtn model =
    Html.div
        [ HtmlAttr.css
            [ Tw.flex
            , Tw.p_2
            ]
        ]
        [ viewAdminModelResult model ]

viewAdminModelResult : AdminModel -> Html.Html Msg
viewAdminModelResult model =
    case model of
        RemoteData.NotAsked ->
            Html.button
                [ rebuildBtnStyle
                    [ Tw.bg_blue_500 ]
                    ( Css.hover [ Tw.bg_blue_700 ] )
                    ( rebuildBtnFocusStyle [ Tw.ring_blue_500, Tw.ring_offset_blue_200 ] )
                , Events.onClick FetchAdminResponse
                ]
                ( rebuildBtnSvg rebuildBtnLogoStyle Icons.build "Rebuild" )

        RemoteData.Loading ->
            Html.button
                [ rebuildBtnStyle
                    [ Tw.bg_yellow_500, Tw.cursor_not_allowed ]
                    ( Css.hover [ Tw.bg_yellow_700 ] )
                    ( rebuildBtnFocusStyle [ Tw.ring_yellow_500, Tw.ring_offset_yellow_200 ] )
                ]
                ( rebuildBtnSvg ( Tw.animate_spin :: rebuildBtnLogoStyle ) Icons.refresh "Building" )

        RemoteData.Failure e ->
            Html.button
                [ rebuildBtnStyle
                    [ Tw.bg_red_500 ]
                    ( Css.hover [ Tw.bg_red_700 ] )
                    ( rebuildBtnFocusStyle [ Tw.ring_red_500, Tw.ring_offset_red_200 ] )
                , Events.onClick FetchAdminResponse
                ]
                ( rebuildBtnSvg rebuildBtnLogoStyle Icons.error_outline ("Error: " ++ Debug.toString e) )

        RemoteData.Success _ ->
            Html.button
                [ rebuildBtnStyle
                    [ Tw.bg_green_500 ]
                    ( Css.hover [ Tw.bg_green_700 ] )
                    ( rebuildBtnFocusStyle [ Tw.ring_green_500, Tw.ring_offset_green_200 ] )
                , Events.onClick FetchAdminResponse
                ]
                (rebuildBtnSvg rebuildBtnLogoStyle Icons.check_circle_outline "Rebuilt")

rebuildBtnStyle : List (Css.Style) -> Css.Style -> Css.Style -> Html.Attribute msg
rebuildBtnStyle cssStyles hoverStyle focusStyle =
    HtmlAttr.css
        ( cssStyles ++
            [ Tw.flex
            , Tw.items_center
            , Tw.shadow
            , Tw.px_4
            , Tw.py_2
            , Tw.text_white
            , Tw.rounded_md
            , hoverStyle
            , focusStyle
            ]
        )

rebuildBtnFocusStyle: List (Css.Style) -> Css.Style
rebuildBtnFocusStyle focusStyles =
    Css.focus
        ( focusStyles ++
            [ Tw.outline_none
            , Tw.ring_2
            , Tw.ring_offset_2
            ]
        )

rebuildBtnSvg : List (Css.Style) -> Icon msg -> String -> List (Html.Html msg)
rebuildBtnSvg cssStyle icon displayText =
    [ Svg.svg
        [ SvgAttr.css cssStyle
        , SvgAttr.viewBox "0 0 24 24"
        ]
        [ Html.fromUnstyled (icon 24 Inherit) ]
    , Html.text displayText
    ]

rebuildBtnLogoStyle : List (Css.Style)
rebuildBtnLogoStyle =
    [ Tw.h_5
    , Tw.w_5
    , Tw.mr_3
    ]

upgradePanel : Html.Html Msg
upgradePanel =
    Html.div
        [ HtmlAttr.css
            [ Tw.text_white
            , Tw.text_base
            , Tw.font_semibold
            , Tw.pt_3
            ]
        ]
        [ Html.a
            [ HtmlAttr.css
                [ Tw.absolute
                , Tw.w_full
                , Tw.bottom_0
                , Tw.bg_blue_700
                , Tw.text_white
                , Tw.flex
                , Tw.items_center
                , Tw.justify_center
                , Tw.py_4
                ]
            ]
            [ Html.text "Upgrade to Pro!" ]
        ]


header : AdminModel -> Html.Html Msg
header model =
    Html.header
        [ HtmlAttr.css
            [ Tw.w_full
            , Tw.items_center
            , Tw.bg_white
            , Tw.py_2
            , Tw.px_6
            , Tw.hidden
            , Bp.sm [ Tw.flex ]
            ]
        ]
        [ rebuildBtn model
        , Html.div
            [ HtmlAttr.css
                [ Tw.w_5over6
                , Tw.justify_center
                , Tw.flex
                , Tw.text_blue_500
                , Tw.text_3xl
                , Tw.font_semibold
                ]
            ]
            [ Html.text "Orbital" ]
        , Html.div
            [ HtmlAttr.css
                [ Tw.relative
                , Tw.w_1over6
                , Tw.flex
                , Tw.justify_end
                ]
            ]
            [ Html.button
                [ HtmlAttr.css
                    [ Tw.relative
                    , Tw.w_12
                    , Tw.h_12
                    , Tw.rounded_full
                    , Tw.overflow_hidden
                    , Tw.border_4
                    , Tw.border_gray_400
                    , Css.hover
                        [ Tw.cursor_default
                        , Tw.border_gray_300
                        ]
                    , Css.focus
                        [ Tw.border_gray_300
                        , Tw.outline_none
                        ]
                    ]
                ]
                [ Html.img [ HtmlAttr.src "https://source.unsplash.com/uJ8LNVCBjFQ/400x400" ][] ]
            ]
        ]

mainSection : List (Html.Html Msg) -> Html.Html Msg
mainSection canvas =
    Html.div
        [ HtmlAttr.css
            [ Tw.bg_gray_200
            , Tw.w_full
            , Tw.overflow_hidden
            , Tw.h_full
            , Tw.border_t
            , Tw.flex
            , Tw.flex_row
            ]
        ]
        [ Html.main_
            [ HtmlAttr.css
                [ Tw.w_full
                , Tw.h_full
                , Tw.flex
                , Tw.flex_grow
                , Tw.pt_6
                , Tw.pr_6
                , Tw.pb_6
                , Tw.flex_row
                , Tw.flex_auto
                ]
            ]
            canvas
        ]

canvasSection : Html.Html Msg -> Html.Html Msg
canvasSection content =
    Html.div
        [ HtmlAttr.css
            [ Tw.bg_white
            , Tw.rounded_md
            , Tw.shadow_xl
            , Tw.flex
            , Tw.flex_auto
            , Tw.flex_wrap
            , Tw.w_2over3
            , Tw.ml_1
            , Tw.p_6
            , Tw.overflow_hidden
            , Bp.lg [ Tw.pr_2 ]
            ]
        ]
        [ Html.div
            [ HtmlAttr.id "graph"
            , HtmlAttr.css
                [ Tw.flex
                , Tw.w_screen
                , Tw.content_center
                , Tw.justify_center
                , Tw.items_center
                , Tw.overflow_hidden
                ]
            ]
            [ content ]
        ]

viewEntityListModelResult : Model -> Html.Html Msg
viewEntityListModelResult model =
    case model.entities of
        RemoteData.NotAsked ->
            Html.text "I didn't ask"

        RemoteData.Loading ->
            Html.div
                [ HtmlAttr.css
                    [ Tw.w_1over3
                    ]
                ]
                [ Html.div
                    [ HtmlAttr.css
                        [ Tw.w_full
                        , Tw.text_center
                        , Tw.justify_center
                        , Tw.text_lg
                        , Tw.inline_flex
                        ]
                    ]
                    [ Html.text "Hold on tight!" ]
                , Html.div
                    [ HtmlAttr.css
                        [ Tw.w_full
                        , Tw.justify_center
                        , Tw.inline_flex
                        , Tw.flex
                        , Tw.pt_8
                        ]
                    ]
                    [ Svg.svg
                        [ SvgAttr.css
                            [ Tw.animate_spin
                            , Tw.h_32
                            , Tw.w_32
                            ]
                        , SvgAttr.viewBox "0 0 24 24"
                        , SvgAttr.fill "none"
                        ]
                        [ Svg.circle
                            [ SvgAttr.css
                                [ Tw.opacity_25 ]
                            , SvgAttr.r "10"
                            , SvgAttr.cx "12"
                            , SvgAttr.cy "12"
                            , SvgAttr.stroke "currentColor"
                            , SvgAttr.strokeWidth "4"
                            ]
                            []
                        , Svg.path
                            [ SvgAttr.css
                                [ Tw.opacity_75 ]
                            , SvgAttr.fill "currentColor"
                            , SvgAttr.d "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            ]
                            []
                        ]
                    ]
                ]

        RemoteData.Failure e ->
            buildFailureMsg e

        RemoteData.Success _ ->
            case model.chartConfig of
                --ChartConfig chartOpts ->
                --    setGraphOptions chartOpts []
                _ ->
                    Html.div
                        [ HtmlAttr.css
                            [ Tw.w_1over3
                            ]
                        ]
                        [ Html.div
                            [ HtmlAttr.css
                                [ Tw.w_full
                                , Tw.text_center
                                , Tw.justify_center
                                , Tw.text_lg
                                , Tw.inline_flex
                                ]
                            ]
                            [ Html.text "I got no options!" ]
                        , Html.div
                            [ HtmlAttr.css
                                [ Tw.w_full
                                , Tw.justify_center
                                , Tw.inline_flex
                                , Tw.flex
                                , Tw.pt_8
                                ]
                            ]
                            [ Html.fromUnstyled (Icons.warning 96 Inherit)
                            ]
                        ]

setGraphOptions : ChartOptions -> List (Html.Html Msg) -> Html.Html Msg
setGraphOptions chartOptions =
    Html.node "echart-element"
    [ HtmlAttr.property "option" <| encodeChartOptions <| chartOptions
    , onNodeClick GotId
    ]

onNodeClick : (Int -> a) -> Html.Attribute a
onNodeClick event =
    Events.on "nodeClick"
        <| Decode.map event detailDataIndexDecoder

detailDataIndexDecoder: Decode.Decoder Int
detailDataIndexDecoder =
    Decode.at [ "detail", "id" ] Decode.int

buildFailureMsg: Error parsedData -> Html.Html Msg
buildFailureMsg parsedData =
    case parsedData of
        Graphql.Http.GraphqlError _ graphqlErrors ->
            buildErrorMsg "Graphql Error" (List.map .message graphqlErrors)

        Graphql.Http.HttpError httpError ->
            buildErrorMsg "Http Error" [ buildHttpErrorMessage httpError ]

buildErrorMsg: String -> List (String) -> Html.Html Msg
buildErrorMsg eType eMsgs =
    Html.div
        [ HtmlAttr.css
            [ Tw.flex
            , Tw.bg_red_200
            , Tw.p_4
            ]
        ]
        [ Html.div
            [ HtmlAttr.css [ Tw.mr_4 ] ]
            [ Html.div
                [ HtmlAttr.css
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
            [ HtmlAttr.css
                [ Tw.flex
                , Tw.justify_between
                , Tw.w_full
                ]
            ]
            [ Html.div
                [ HtmlAttr.css [ Tw.text_red_600 ] ]
                [ Html.p
                    [ HtmlAttr.css
                        [ Tw.mb_2
                        , Tw.font_bold
                        , Tw.text_lg
                        ]
                    ]
                    [ Html.text eType ]
                , Html.ul
                    [ HtmlAttr.css
                        [ Tw.text_base ]
                    ]
                    (List.map (\msg -> Html.li [][ Html.text msg ]) eMsgs)
                ]
            , Html.div
                [ HtmlAttr.css
                    [ Tw.text_sm
                    , Tw.text_gray_500
                    ]
                ]
                [ Html.button [] [ Html.fromUnstyled (Icons.refresh 24 Inherit) ] ]
            ]
        ]

buildHttpErrorMessage : HttpError -> String
buildHttpErrorMessage httpError =
    case httpError of
        Graphql.Http.BadUrl message ->
            message

        Graphql.Http.Timeout ->
            "Server is taking too long to respond. Please try again later."

        Graphql.Http.NetworkError ->
            "Unable to reach server."

        Graphql.Http.BadStatus metadata body ->
            "Request failed with status code: " ++ String.fromInt metadata.statusCode ++ ". Error: " ++ body

        Graphql.Http.BadPayload error ->
            "Bad payload received: " ++ errorToString error

entityDetails : SelectedEntityModel -> Html.Html Msg
entityDetails model =
    case model of
        Selected entity ->
            Html.div
                [ HtmlAttr.css
                    [ Tw.relative
                    , Tw.flex
                    , Tw.flex_auto
                    , Tw.w_full
                    ]
                ]
                [ Html.div
                    [ HtmlAttr.css
                        [ Tw.relative
                        , Tw.w_screen
                        , Tw.max_w_md
                        ]
                    ]
                    [ Html.div
                        [ HtmlAttr.css
                            [ Tw.h_full
                            , Tw.flex
                            , Tw.flex_col
                            , Tw.bg_white
                            , Tw.overflow_y_hidden
                            ]
                        ]
                        [ Html.dl
                            []
                            [ Html.div (rowStyle Tw.bg_gray_50) (rowContent "name" [ Html.text entity.name ])
                            , Html.div (rowStyle Tw.bg_white) (rowContent "context" [ Html.text entity.context ])
                            , Html.div (rowStyle Tw.bg_gray_50) (rowContent "graph-id" [ Html.text <| toString entity.id ])
                            , Html.div
                                (rowStyle Tw.bg_white)
                                (rowContent "connections"
                                    <| List.concatMap (\v -> [ Html.text v, Html.br [][] ]) entity.connections
                                )
                            ]
                        ]
                    ]
                ]
        NoEnt ->
            Html.div []
                [ Html.div
                    [ HtmlAttr.css
                        [ Tw.justify_center
                        , Tw.inline_flex
                        , Tw.flex
                        , Tw.w_full
                        , Tw.pt_36
                        ]
                    ]
                    [ Html.div
                        [ HtmlAttr.css
                            [ Tw.bg_gray_100
                            , Tw.rounded_full
                            , Tw.flex
                            , Tw.p_4
                            ]
                        ]
                        [ Svg.svg
                            [ SvgAttr.css
                                [ Tw.h_32
                                , Tw.w_32
                                ]
                            , SvgAttr.viewBox "0 0 24 24"
                            ]
                            [ Html.fromUnstyled (Icons.travel_explore 24 Inherit) ]
                        ]
                    ]
                , Html.div
                    [ HtmlAttr.css
                        [ Tw.pt_8
                        , Tw.justify_center
                        , Tw.w_full
                        , Tw.inline_flex
                        , Tw.font_bold
                        , Tw.text_lg
                        ]
                    ]
                    [ Html.text "No Entity Found" ]
                , Html.div
                    [ HtmlAttr.css
                        [ Tw.pt_2
                        , Tw.px_8
                        , Tw.text_gray_700
                        , Tw.text_base
                        , Tw.justify_center
                        , Tw.text_center
                        , Tw.w_full
                        , Tw.inline_flex
                        ]
                    ]
                    [ Html.text "No worries! Select a node from the graph to see its properties. You can also filter nodes from the legend." ]
                ]


rowStyle : Css.Style -> List (Html.Attribute Msg)
rowStyle bgColor =
    [ HtmlAttr.css
        [ bgColor
        , Tw.px_4
        , Tw.py_5
        , Bp.sm
            [ Tw.grid
            , Tw.grid_cols_3
            , Tw.gap_4
            , Tw.px_6
            ]
        ]
    ]

rowContent : String -> List (Html.Html Msg) -> List (Html.Html Msg)
rowContent key value =
    [ Html.dt
        [ HtmlAttr.css
            [ Tw.text_sm
            , Tw.font_bold
            , Tw.text_gray_500
            ]
        ]
        [ Html.text key ]
    , Html.dd
        [ HtmlAttr.css
            [ Tw.mt_1
            , Tw.text_sm
            , Tw.text_gray_900
            , Bp.sm
                [ Tw.mt_0
                , Tw.col_span_2
                ]
            ]
        ]
        value
    ]
