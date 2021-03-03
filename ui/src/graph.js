import '@webcomponents/custom-elements'
import * as echarts from 'echarts/index.blank';
import 'echarts/lib/component/title';
import 'echarts/lib/component/tooltip';
import 'echarts/lib/component/legend';
import 'echarts/lib/chart/graph';
import 'zrender/lib/canvas/canvas';

customElements.define('echart-element',
  class EChartElement extends HTMLElement {

    constructor() {
      super();
      this.chart = null;
      this._option = null;
    }

    connectedCallback () {
      const option = this.option
      this.chart = echarts.init(document.getElementById('graph'));
      this.chart.setOption(option);
    }

    set option (newValue) {
      this._option = newValue;
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
