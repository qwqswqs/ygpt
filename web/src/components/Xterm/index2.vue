<template>
  <div id="k8s-xterm" ref="xtermWrapper" class="mt-3">
    <div id="xterm" style="min-height: 800px;" />
  </div>
</template>

<script>
import querystring from 'qs'
import { Terminal } from 'xterm'
import * as fit from 'xterm/lib/addons/fit/fit'
import 'xterm/src/xterm.css'
import io from 'socket.io-client'

export default {
  name: 'Xterm',
  props: {
    connectParams: {
      type: String,
      required: true,
    },
    height: {
      type: Number,
      default: 800,
    },
  },
  data () {
    return {
      socket: null,
      term: null,
    }
  },
  mounted() {
    this.checkAndInitTerminal();
  },
  watch: {
    connectParams: {
      handler (val, oval) {
        if (val) {
          this.checkAndInitTerminal()
        } else {
          this._resetDom()
        }
      },
    },
    height (val, oval) {
      if (val) {
        const terminalDom = document.getElementById('xterm')
        terminalDom.style.minHeight = val + 'px'
        this.term && this.term.fit()
      }
    },
  },
  beforeDestroy () {
    this._socketClose()()
  },
  methods: {
    checkAndInitTerminal() {
      if (this.connectParams) {
        this.initTerminal();
      } else {
        this._resetDom();
        // 可能还需要其他清理工作
      }
    },
    _resetDom() {
      const wrapper = this.$refs.xtermWrapper
      console.log(wrapper)
      const terminalDom = document.getElementById('xterm')
      if (terminalDom && wrapper) {
        wrapper.removeChild(terminalDom)
      }
    },
    _createDom() {
      this._resetDom()
      const wrapper = this.$refs.xtermWrapper
      const newTerminalDom = document.createElement('div')
      newTerminalDom.setAttribute('id', 'xterm')
      wrapper.appendChild(newTerminalDom)
      newTerminalDom.style.minHeight = this.height + 'px'
      return newTerminalDom
    },
    initTerminal() {
      if (!this.connectParams) {
        console.warn('connectParams is not provided, skipping terminal initialization.');
        return;
      }
      // console.log(this.connectParams)
      // const param = querystring.parse(this.connectParams)
      // console.log(param)
      Terminal.applyAddon(fit)
      this.socket = io(this.connectParams)
      this.term = new Terminal({
        cols: 80,
        rows: 24,
        ScrollBar: true,
      })
      const terminalDom = this._createDom()
      this.term.open(terminalDom)
      this.term.focus()
      // term.fit()
      this.term.on('resize', size => {
        this.socket.emit('resize', [size.cols, size.rows])
      })
      this.term.on('data', data => this.socket.emit('input', data))
      this.term.on('keypress', (val, e) => {
        e.preventDefault()
      })
      this.socket.on('output', arrayBuffer => {
        this.term.write(arrayBuffer)
      })
      this.socket.on('disconnection', this._socketClose())
      this.socket.on('connect_error', this._socketClose())
      this.socket.on('disconnect', this._socketClose())
      this.socket.on('disconnecting', this._socketClose())
      window.addEventListener('resize', () => {
        this.term.fit()
      })
      this.term.fit()
    },
    _socketClose() {
      return () => {
        console.log('Connection lose!!!')
        this.socket && this.socket.close()
        this.$emit('close')
      }
    },
  },
}
</script>
