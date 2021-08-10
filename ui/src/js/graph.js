import "@webcomponents/custom-elements";
import * as echarts from "echarts/index.blank";
import "echarts/lib/component/title";
import "echarts/lib/component/tooltip";
import "echarts/lib/component/legend";
import "echarts/lib/chart/graph";
import "zrender/lib/canvas/canvas";

customElements.define(
  "echart-element",
  class EChartElement extends HTMLElement {
    constructor() {
      super();
      this._option = null;
    }

    render(option) {
      const dom = document.getElementById("graph");
      const instance = echarts.getInstanceByDom(dom);
      const elem = this;

      if (instance) {
        echarts.dispose(instance);
      }
      const chart = echarts.init(dom);
      chart.setOption(option);
      chart.on("click", { dataType: "node" }, function (params) {
        console.log("You clicked " + params.name + ", index: " + params.dataIndex);

        elem.dispatchEvent(
          new CustomEvent("nodeClick", {
            bubbles: false,
            detail: {
              id: params.dataIndex,
            },
          }),
        );
      });
    }

    connectedCallback() {
      this.render(this.option);
    }

    set option(newValue) {
      this._option = newValue;
    }

    get option() {
      return this._option;
    }

    disconnectedCallback() {}
  },
);
