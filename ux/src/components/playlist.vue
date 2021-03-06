<template>
  <div class="playlist">
    <section class="play">
      <h1>Listen</h1>
      <audio xmlns="http://www.w3.org/1999/xhtml" controls="controls" preload="none" style="width: 100%;">
        <source type="application/ogg" :src="stream">
      </audio>
    </section>
    <section class="add">
      <h1>Add a new track, from YouTube</h1>
      <form @submit.prevent="upload()">
        <input type="text" v-model="url" placeholder="Enter a YouTube URL here..." />
        <button>&plus;</button>
      </form>
    </section>
    <section class="current">
      <h1>Current Track Line-up</h1>
      <ul>
        <li v-for="song in songs" :key="song.file">
          <input type="checkbox" :checked="song.active == 1"
                 @change="song.active = song.active ? 0 : 1; dirty()" />
          <label>{{ pretty(song.file) }}</label>
        </li>
      </ul>
      <button @click="save()" :class="{saving, saved}">
        <template v-if="saving">Saving...</template>
        <template v-if="saved">Saved!</template>
        <template v-else>Save Changes</template>
      </button>
    </section>
  </div>
</template>

<script>
export default {
  props: {
    stream: String
  },
  data() {
    return {
      url:    '',
      timer:  null,
      saving: false,
      saved:  false,
      songs: [],
    }
  },
  methods: {
    pretty(file) {
      return file.replace(/^\/radio\//, '').replace(/-*-[A-Za-z0-9]+\..*?\.opus/, '')
    },

    sync() {
      fetch('/playlist')
        .then(r => r.json())
        .then(that => this.songs = that.playlist)
    },

    dirty() {
      this.saved = false
      window.clearTimeout(this.timer)
    },

    save() {
      fetch('/playlist', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ playlist: this.songs }),
      })
        .then(r => r.json())
        .then(() => {
          this.saving = false;
          this.saved = true;
          this.timer = window.setTimeout(() => this.saved = false, 1800);
          this.sync()
        })
    },

    upload() {
      fetch('/upload', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: this.url }),
      })
        .then(r => r.json())
        .then(() => this.url = '')
    }
  },
  mounted() {
    this.sync()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
  margin: 0 0 0.5em 0;
}
section + section {
  margin-top: 50pt;
}
section + section h1 {
  border-bottom: 1px dashed #777;
}
ul {
  list-style-type: none;
  padding: 0;
  text-align: left;
}
li {
  display: block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
form {
  display: flex;
  flex-direction: row;
}
form input[type=text] {
  font-size: 18pt;
  width: 100%;
  padding: 0.5em;
  box-sizing: border-box;
}
button {
  cursor: pointer;
  background-color: #a43ed2;
  font-weight: bold;
  color: #fff;
  border: 1px solid #ccc;
  font-size: 14pt;
  padding: 9pt 12pt;
}
form button {
}
form input,
form button {
  margin: 0;
  border: 1px solid #ccc;
}
.current button {
  font-size: 12pt;
}
.current button.saving {
  background-color: #73617b;
}
.current button.saved {
  background-color: #12a652;
}
</style>
