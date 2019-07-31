var base_url = "http://localhost:8080";

(function(){
  if (Cookies.get("login") == null) {
    window.location.href = base_url+"/static/login.html"
  }
})()

Vue.component('timer', {
  data: function() {
    return {
      baseUrl: "http://localhost:8080",

      timeBegan : null,
      timeStopped : null,
      stoppedDuration : 0,
      started : null,
      running : false,
      time: '00:00:00.000',

      recordTitle: '',
      disabled: false
    }
  },
  methods: {
    start: function() {
      if (this.running) return
      if (this.timeBegan === null) {
        this.reset()
        this.timeBegan = moment.utc()
      }
      if (this.timeStopped !== null) {
        this.stoppedDuration += (moment.utc() - this.timeStopped)
      }
      this.started = setInterval(this.clockRunning, 10)
      this.running = true
    },
    stop: function() {
      this.running = false
      this.timeStopped = moment.utc()
      clearInterval(this.started)
    },
    reset: function() {
      this.running = false
      clearInterval(this.started)
      this.stoppedDuration = 0
      this.timeBegan = null
      this.timeStopped = null
      this.time = "00:00:00.000"
    },
    clockRunning: function() {
      var te = moment.utc(moment.utc() - this.timeBegan - this.stoppedDuration)
      , hour = te.hours()
      , min = te.minutes()
      , sec = te.seconds()
      , ms = te.milliseconds()

      this.time =
        this.zeroPrefix(hour, 2) + ":" +
        this.zeroPrefix(min, 2) + ":" +
        this.zeroPrefix(sec, 2) + "." +
        this.zeroPrefix(ms, 3)
    },
    zeroPrefix: function(num, digit) {
      var zero = ''
      for(var i = 0; i < digit; i++) {
        zero += '0'
      }
      return (zero + num).slice(-digit)
    },
    sendRecord: function() {
      data = {
        time: this.time,
        title: this.recordTitle,
        method: 'add'
      }
      $.post({
        url: base_url+'/v1/timer',
        data: data,
        success: function(res){
          console.log(res);
        }
      })
      this.disabled = true
    }
  },
  template: `
    <div class="row">
      <div class="col-md-4">
        <div class="input-group input-group-lg mb-3">
          <div class="input-group-prepend">
            <button v-bind:disabled="disabled" v-on:click="start" class="btn btn-outline-secondary" type="button" title="start">start</button>
          </div>
          <div class="btn btn-light d-flex align-items-center">{{time}}</div>
          <div class="input-group-append">
            <button v-bind:disabled="disabled" v-on:click="stop" v-on:dblclick="reset" class="btn btn-outline-secondary" type="button" title="double click to reset">stop/reset</button>
          </div>
        </div>
      </div>
      <div class="col">
        <div class="input-group input-group-lg mb-3">
          <div class="input-group-prepend">
            <button v-on:click="sendRecord" v-bind:disabled="disabled" v-bind:class="[disabled ? 'btn-success' : 'btn-outline-success']" class="btn" type="button" title="save to database">&#10004 <span v-if="disabled">saved</span></button>
          </div>
          <input v-model="recordTitle" v-bind:readonly="disabled" class="form-control col-md-7" type="text" placeholder="Set Title">
          <div class="input-group-append">
            <button v-on:click="$emit('removeRecord')" class="btn btn-outline-danger" type="button" title="remove record">&#10006</button>
          </div>
        </div>
      </div>
    </div>
  `
})

Vue.component('custom-table', {
  data: function() {
    return {
      timer: ['']
    }
  },
  methods: {
    addTimer: function() {
      this.timer.push('')
    }
  },
  template: `
    <div class="p-5">
      <button v-on:click="addTimer" class="btn btn-success">add</button>
      <table class="table table-borderless my-3">
        <tbody>
          <tr v-for="(v, i) in timer">
            <td><timer v-on:removeRecord="timer.splice(i, 1)"></timer></td>
          </tr>
        </tbody>
      </table>
    </div>
  `
})

new Vue({
  el: '#app'
})
