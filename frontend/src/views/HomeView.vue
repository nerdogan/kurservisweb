<style>
@media (min-width: 512px) {
  .about {
    min-height: 120vh;
    display: flex;
    align-items: center;
  }
.button-21 {
  align-items: center;
  appearance: none;
  background-color: #3EB2FD;
  background-image: linear-gradient(1deg, #4F58FD, #149BF3 99%);
  background-size: calc(100% + 20px) calc(100% + 20px);
  border-radius: 100px;
  border-width: 0;
  box-shadow: none;
  box-sizing: border-box;
  color: #FFFFFF;
  cursor: pointer;
  display: inline-flex;
  font-family: CircularStd,sans-serif;
  font-size: 1rem;
  height: auto;
  justify-content: center;
  line-height: 1.5;
  padding: 6px 20px;
  position: relative;
  text-align: center;
  text-decoration: none;
  transition: background-color .2s,background-position .2s;
  user-select: none;
  -webkit-user-select: none;
  touch-action: manipulation;
  vertical-align: top;
  white-space: nowrap;
}

.button-21:active,
.button-21:focus {
  outline: none;
}

.button-21:hover {
  background-position: -20px -20px;
}

.button-21:focus:not(:active) {
  box-shadow: rgba(40, 170, 255, 0.25) 0 0 0 .125em;
}
input[type=text] {
  width: 50%;
  padding: 3px 5px;
  margin: 4px 0;
  font-size: 14px;
  box-sizing: border-box;
  border: 3px solid #ccc;
  -webkit-transition: 0.5s;
  transition: 0.5s;
  outline: none;
}

input[type=text]:focus {
  border: 3px solid #555;
}
}
</style>
<template>
  
    <div>
    
    <br>
    <form @submit.prevent="handleSubmit">
 
      <p>Kaç Gram ?</p>  
  <input type="text" inputmode="decimal"  v-model="gram">
 
  <p>Kaç Milyem ?</p>  
  <input type="text" inputmode="decimal" v-model="milyem">
 
  <p>Altın Dolar Fiyatı ?</p>  
  <input type="text" inputmode="decimal" v-model="usd" >
  <br><br>

<button class="button-21" role="button" type="submit">Hesapla</button>

  <br><br>
<p>Sonuç: 
  <h1> {{ description }}  USD</h1>
  
</p>
    </form>
    <p>
      <h1> &nbsp</h1>

  {{ tarih }}
  <p>
  Namık ERDOĞAN © 2024
</p>
</p>

</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import axios from 'axios'

interface Item {
tarih: any
  tutar: number
  masano: string
}
export default defineComponent({
  name: 'ItemList112',
  data() {
    return {
      items: [] as Item[],
      usd:0.0,
      description: 0.0,
      gram:"1",
      gramnum:0.0 ,
      milyem: null ,
      tarih:" "
    }
  },
  mounted() {
    axios.get<Item[]>('http://localhost:8080/price?productId=3')
      .then(response => {
        console.log(response.data)
        this.items = response.data
        this.usd= (response.data[0].tutar)*1
        this.tarih=response.data[0].tarih
      })
      .catch(error => {
        console.error(error)
      })
  },
  methods: {
    handleSubmit() {
      this.gramnum=parseFloat(this.gram.replace(",","."))
      this.description=parseFloat(((this.usd+100.0)*this.milyem*this.gramnum/995000).toFixed(2))
    }
  }
})
</script>
