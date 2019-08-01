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

Vue.component('layout', {
  template: `
    <div class="p-5">
      <a href="index.html" class="btn btn-info float-right">back</a>
      <nav>
        <div class="nav nav-tabs" id="nav-tab" role="tablist">
          <a class="nav-item nav-link active" id="nav-home-tab" data-toggle="tab" href="#nav-home" role="tab" aria-controls="nav-home" aria-selected="true">Home</a>
          <a class="nav-item nav-link" id="nav-history-tab" data-toggle="tab" href="#nav-history" role="tab" aria-controls="nav-history" aria-selected="false">History</a>
        </div>
      </nav>
      <div class="tab-content" id="nav-tabContent">
        <div class="tab-pane fade show active" id="nav-home" role="tabpanel" aria-labelledby="nav-home-tab">
          <tab-one></tab-one>
        </div>
        <div class="tab-pane fade" id="nav-history" role="tabpanel" aria-labelledby="nav-history-tab">
          <tab-two></tab-two>
        </div>
      </div>
    </div>
  `
})

Vue.component('tab-one', {
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
    <div class="p-2">
      <button v-on:click="addTimer" class="btn btn-success">add</button>
      <table id="table-timer" class="table table-borderless my-3">
        <tbody>
          <tr v-for="(v, i) in timer">
            <td><timer v-on:removeRecord="timer.splice(i, 1)"></timer></td>
          </tr>
        </tbody>
      </table>
    </div>
  `
})

Vue.component('tab-two', {
  mounted: function() {
    $('#table-history').DataTable({
      "processing": true,
      "ajax": {
        "type": "POST",
        "url": base_url+"/v1/timer",
        "data": {method:"get"},
        "dataSrc": function (src) {
          var data = new Array(), ctr = 1;
          for (var i = 0; i < src.data.length; i++) {
            data.push({
              "No" : ctr++,
              "Title" : src.data[i].Title,
              "Time" : src.data[i].Time
            })
          }
          return data
        }
      },
      "columns" : [
        { "data": "No", "width": "20%" },
        { "data": "Title" },
        { "data": "Time" }
      ]
    })
  },
  template: `
    <div class="p-3">
      <table id="table-history" class="display table table-striped table" style="width:100%">
        <thead>
            <tr>
              <th>No</th>
              <th>Title</th>
              <th>Time</th>
            </tr>
        </thead>
      </table>
    </div>
  `
})

new Vue({
  el: '#app'
})
