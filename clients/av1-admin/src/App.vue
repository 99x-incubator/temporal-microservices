<template>
  <div class="container mt-5">
    <h2 class="text-center mb-4">Robot Control Panel</h2>
    <div class="card p-3 shadow-sm">
      <ul class="list-group">
        <li
          v-for="robot in robots"
          :key="robot.id"
          class="list-group-item d-flex justify-content-between align-items-center"
        >
          <span>
            ({{ robot.id }}) {{ robot.name }}  <small v-if="robot.disabled" class="text-danger">(Disabled)</small>
          </span>
          <button
            v-if="!robot.disabled"
            @click="disableRobot(robot.id)"
            class="btn btn-danger btn-sm"
            :disabled="robot.disabled"
          >
            Disable
          </button>
        </li>
      </ul>
    </div>
    <div
      v-if="notification"
      variant="info"
      dismissible
      class="mt-3"
      @dismissed="notification = ''"
      style="background-color: #d1ecf1; color: #0c5460; border-color: #bee5eb;"
    >
      <div class="text-right mt-2">
        {{ notification }}
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      userID: "user123", // Example userID, could be dynamic
      robots: [
        { id: 1, name: 'Robot Alpha', disabled: false },
        { id: 2, name: 'Robot Beta', disabled: false },
        { id: 3, name: 'Robot Gamma', disabled: false },
      ],
      notification: '',
    };
  },
  mounted() {
    // Establish a WebSocket connection with the userID as a query parameter
    var loc = window.location
    const socket = new WebSocket(`ws://` + loc.host + `/notification/ws?userID=${this.userID}`);

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("WebSocket Message:", data);

      this.notification = data.message;
    };

    socket.onerror = (error) => {
      console.error("WebSocket Error:", error);
    };
  },
  methods: {
    async disableRobot(robotId) {
      this.notification = ''; // Clear previous notification
      
      // Call the /disable_robot endpoint of the gateway microservice
      const response = await axios.post('api/disable_robot', {
        robot_id: robotId.toString(),
        user_id: this.userID,
      });

      if (response.status === 200) {
        this.notification = 'Your request submitted, please wait.';
      }

      const robot = this.robots.find((r) => r.id === robotId);
      robot.disabled = true;
    },
  },
};
</script>

<style>
body {
  background-color: #f4f4f9;
}

.container {
  max-width: 600px;
}

h2 {
  color: #343a40;
}

.card {
  background-color: #ffffff;
  border-radius: 8px;
}

.btn-danger {
  background-color: #ff5c5c;
  border: none;
}

.btn-danger:hover {
  background-color: #ff4242;
}
</style>
