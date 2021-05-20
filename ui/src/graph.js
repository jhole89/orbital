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
            this._option = null;
            console.log("contructor called")
        }

        render(option) {
          const dom = document.getElementById('graph');
          const instance = echarts.getInstanceByDom(dom);

          console.log("fetched instance: " + instance)
          if (instance) {
            instance.showLoading();
            console.log("disposing...instance " + instance)
            echarts.dispose(instance);
          }
          if (!instance) {

          }
          const chart = echarts.init(dom);
          console.log("chart init'd..." + chart)
          chart.setOption(option);
        }

        connectedCallback () {
          this.render(this.option);
        }

        set option (newValue) {
          this._option = newValue;
        }

        get option () {
            console.log("get option")
            return this._option
        }

        disconnectedCallback () {}
    }
);
