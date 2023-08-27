import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    isConnected: false,
    connectMethod: 0, // 1: walletconnect, 2: metamask
    address: "",
    chainID: 0,

    ethChainAddress: "",
    ethContractAddress: "",
    ethLibraries: ["", ""],

    cashboxChainAddress: "",
    cashboxDataContractAddress: "",
    cashboxControllerContractAddress: "",
    cashboxLibraries: ["", ""]
  },
})
