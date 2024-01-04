async function deleteUser(key, event) {
    await axios.delete(`https://webprojass4-default-rtdb.firebaseio.com/users/${key}.json`)
}