<template>
  <section class="hero is-primary is-fullheight">
    <div class="hero-body">
      <div class="container">
        <div class="columns is-centered">
          <div class="column is-5-tablet is-4-desktop is-3-widescreen">
            <b-notification
              v-if="login_error"
              type="is-danger"
              role="alert"
              :closable="closable"
            >Wrong email/password</b-notification>

            <div class="box">
              <div class="field">
                <label for class="label">Email</label>
                <div class="control has-icons-left">
                  <input
                    type="email"
                    placeholder="e.g. bobsmith@gmail.com"
                    class="input"
                    required
                    v-model="email"
                  />
                  <span class="icon is-small is-left">
                    <i class="fa fa-envelope"></i>
                  </span>
                </div>
              </div>
              <div class="field">
                <label for class="label">Password</label>
                <div class="control has-icons-left">
                  <input
                    type="password"
                    placeholder="*******"
                    class="input"
                    required
                    v-model="password"
                  />
                  <span class="icon is-small is-left">
                    <i class="fa fa-lock"></i>
                  </span>
                </div>
              </div>
              <div class="field">
                <button class="button is-success" @click="login">Login</button>
              </div>
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
      email: "",
      password: "",
      closable: false,
      login_error: false
    };
  },
  created() {
    const token = localStorage.getItem("token");
    if (token != "") {
      this.$router.push("chat");
    }
  },
  methods: {
    login() {
      this.login_error = false;
      this.$http
        .post("http://127.0.0.1:9001/sign-in", {
          email: this.email,
          password: this.password
        })
        .then(response => {
          localStorage.setItem("token", response.data.token);
          this.$router.push("chat");
        })
        .catch(error => {
          this.login_error = true;
          // eslint-disable-next-line no-console
          console.log(error);
        });
    }
  }
};
</script>