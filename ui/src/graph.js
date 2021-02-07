import '@webcomponents/custom-elements'
import * as echarts from 'echarts/index.blank';
import 'echarts/lib/component/title';
import 'echarts/lib/component/tooltip';
import 'echarts/lib/component/legend';
import 'echarts/lib/chart/graph';
import 'zrender/lib/canvas/canvas';
import 'echarts/lib/chart/bar';

customElements.define('echart-element',
  class EChartElement extends HTMLElement {
    constructor() {
      super();
      this.chart = null;
      this._option = null;
    }
    connectedCallback () {
      this.chart = echarts.init(document.getElementById('graph'));
      console.log("CONNECTED CALLBACK: this.chart is: ");
      console.log(this.chart);
    }

    set option (newValue) {
      console.log("SET OPTION: newValue is: ");
      console.log(newValue);
      this._option = newValue;
      console.log("SET OPTION: this.chart is: ")
      console.log(this.chart)
      if (this.chart) {
        this.chart.setOption(newValue);
      }
    }

    get option () {
      return this._option
    }

    disconnectedCallback () {}
  }
);