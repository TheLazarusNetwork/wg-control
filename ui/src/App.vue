<template>
  <v-app id="inspire">

    <v-app-bar app>
      <img class="mr-3" :src="require('./assets/logo.png')" height="50" alt="WireGuard Control"/>
      <v-toolbar-title to="/">WireGuard Control</v-toolbar-title>

      <v-spacer />

      <v-toolbar-items>
        <v-btn to="/clients">
          Clients
          <v-icon right dark>mdi-account-network-outline</v-icon>
        </v-btn>
        <v-btn to="/server">
          Server
          <v-icon right dark>mdi-vpn</v-icon>
        </v-btn>
      </v-toolbar-items>

    </v-app-bar>

    <v-content>
      <v-container>
        <router-view />
      </v-container>
      <Notification v-bind:notification="notification"/>
    </v-content>

    <v-footer app>
      <v-row justify="start" no-gutters>
        <v-col cols="12" lg="6" md="12" sm="12">
          <div :align="$vuetify.breakpoint.smAndDown ? 'center' : 'left'">
            <small>Copyright &copy {{ new Date().getFullYear() }} | Lazarus Network</small>
          </div>
        </v-col>
      </v-row>
      <v-row justify="end" no-gutters>
        <v-col cols="12" lg="6" md="12" sm="12">
          <div :align="$vuetify.breakpoint.smAndDown ? 'center' : 'right'">
            <small>Created with</small>
            <v-icon class="pr-1 pl-1">mdi-heart</v-icon><span>using</span><a class="pr-2 pl-1" href="mailto:connect@lazarus.network">GoLang & Vue.js</a>
            <a :href="'https://github.com/TheLazarusNetwork/'"><kbd>Version: {{ version.substring(0,7) }}</kbd></a>
          </div>
        </v-col>
      </v-row>
    </v-footer>

  </v-app>
</template>

<script>
  import {ApiService} from "./services/ApiService";
  import Notification from './components/Notification'

  export default {
    name: 'App',

    components: {
      Notification
    },

    data: () => ({
      api: null,
      version:'N/A',
      notification: {
        show: false,
        color: '',
        text: '',
      },
    }),

    mounted() {
      this.api = new ApiService();
      this.getVersion()
    },

    created () {
      this.$vuetify.theme.dark = true
    },

    methods: {
      getVersion() {
        this.api.get('/server/version').then((res) => {
          this.version = res.version;
        }).catch((e) => {
          this.notify('error', e.response.status + ' ' + e.response.statusText);
        });
      },
      notify(color, msg) {
        this.notification.show = true;
        this.notification.color = color;
        this.notification.text = msg;
      }
    }
  };
</script>
