<template>
  <div>
    <xterm :connectParams="connectParams" />
  </div>
</template>
<script>
import xterm from "../../../../src/components/Xterm/index2.vue"
export default {
  components:{
    xterm,
  },
  data() {
    return {
      connectParams: '',
    };
  },
  created() {
    this.connectParams=this.$route.query.id
  }
}

</script>