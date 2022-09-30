import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    address: "",
    chainID: 0,
    reConnect: 0,

    ethChainAddress: "",
    ethContractAddress: "",
    ethLibraries: ["", "", "", ""],

    cashboxChainAddress: "",
    cashboxDataContractAddress: "",
    cashboxControllerContractAddress: "",
    cashboxLibraries: ["", "", "", "", ""]
  },
})
