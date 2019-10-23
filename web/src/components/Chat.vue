<template>
  <section class="hero is-primary is-fullheight">
    <div class="hero-body">
      <div class="container">
        <div class="columns is-centered">
          <div class="box" style="width:600px; height: 800px;position: relative;">
            <!-- Display messages -->
            <div style="overflow: auto; height: 620px;margin-botom: 20px;">
              <template v-for="(message, id) in messages">
                <div
                  class="notification"
                  v-bind:key="id"
                  v-bind:class="{ 'is-success': message.to }"
                >
                  <strong>{{ message.from }}</strong>
                  <template v-if="message.to">
                    <small>&nbsp;private</small>
                  </template>
                  : {{ message.text }}
                </div>
              </template>
            </div>

            <!-- Display controls -->
            <b-field grouped style="margin-top: 10px;">
              <b-field label="Name">
                <b-input v-model="messageReceiver"></b-input>
              </b-field>
              <b-field label="Message" expanded>
                <b-input v-model="messageInput"></b-input>
              </b-field>
            </b-field>

            <div class="buttons">
              <button class="button is-primary" @click="send">Send</button>
              <button class="button is-danger" @click="logout">Logout</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      messageInput: "",
      messageReceiver: "",
      messages: [],
      ws: null
    };
  },
  created() {
    // We should connect websocket here
    const token = localStorage.getItem("token");
    if (token == "") {
      this.$router.push("/");
    }
    const ws = new WebSocket("ws://127.0.0.1:9002?token=" + token);

    ws.onmessage = message => {
      try {
        this.messages.push(JSON.parse(message.data));
      } catch (e) {
        this.error = e;
      }
    };
    ws.onerror = error => {
      this.error = error;
    };
    ws.onopen = () => {
      this.ws = ws;
    };
  },
  methods: {
    logout() {
      localStorage.setItem("token", "");
      this.$router.push("/");
    },
    send() {
      var message = {
        to: this.messageReceiver,
        text: this.messageInput
      };

      this.messageInput = "";
      this.messageReceiver = "";
      this.ws.send(JSON.stringify(message));
    }
  }
};
</script>

<style scoped>
.input {
  margin-bottom: 10px;
}
</style>