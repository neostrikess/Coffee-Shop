package main

import (
	"coffee-shop-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// OrderStatus représente l'état d'une commande

// Base de données en mémoire
var drinks []models.Drink
var orders []models.Order
var orderCounter int = 1

func getMenu(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type à application/json
	w.Header().Set("Content-Type", "application/json")
	// 2. Encoder et retourner le slice drinks en JSON
	json.NewEncoder(w).Encode(drinks)
}

func getDrink(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Récupérer l'ID depuis les variables de route (mux.Vars)*
	vars := mux.Vars(r)
	id := vars["id"]

	// 3. Parcourir le slice drinks
	for _, drink := range drinks {
		// 4. Si l'ID correspond à celui recherché
		if drink.ID == id {
			// 5. Encoder et retourner la boisson en JSON
			json.NewEncoder(w).Encode(drink)
			return
		}
	}
	// 6. Sinon : retourner une erreur 404
	http.Error(w, "Boisson non trouvée", http.StatusNotFound)

}
func createOrder(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Décoder le body JSON dans une variable Order
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// 3. Gérer l'erreur de décodage (400 Bad Request)
	// (déjà fait ci-dessus)
	// 4. Vérifier que la boisson (DrinkID) existe dans drinks
	var drinkFound *models.Drink
	for _, drink := range drinks {
		if drink.ID == order.DrinkID {
			drinkFound = &drink
			break
		}
	}
	if drinkFound == nil {
		http.Error(w, "Boisson non trouvée", http.StatusBadRequest)
		return
	}

	// 5. Si non trouvée : retourner 400 Bad Request
	// 6. Générer un ID unique (ex: ORD-001, ORD-002...)
	order.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	// 7. Remplir order.DrinkName avec le nom de la boisson
	order.DrinkName = drinkFound.Name

	// 8. Définir order.Status à StatusPending
	order.Status = models.StatusPending

	// 9. Définir order.OrderedAt à time.Now()
	order.OrderedAt = time.Now()

	// 10. Calculer le prix total (appeler calculatePrice)
	order.TotalPrice = calculatePrice(drinkFound.BasePrice, order.Size, order.Extras)

	// 11. Ajouter la commande au slice orders
	orders = append(orders, order)

	// 12. Retourner 201 Created avec la commande en JSON
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)

}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func calculatePrice(basePrice float64, size string, extras []string) float64 {
	// TODO:
	// 1. Partir du basePrice
	price := basePrice

	// 2. Ajuster selon la taille:
	//    - "small" : x0.8
	//    - "medium" : x1.0
	//    - "large" : x1.3
	switch size {
	case "small":
		price *= 0.8
	case "medium":
		price *= 1.0
	case "large":
		price *= 1.3
	default:
		// Taille inconnue, on peut choisir de ne rien faire ou retourner une erreur
	}
	// 3. Ajouter 0.50€ pour chaque extra
	price += float64(len(extras)) * 0.50

	// 4. Retourner le prix total
	return price

}
func getOrders(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")

	// 2. Encoder et retourner le slice orders en JSON
	json.NewEncoder(w).Encode(orders)
}
func getOrder(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Récupérer l'ID depuis les variables de route
	vars := mux.Vars(r)
	id := vars["id"]

	// 3. Parcourir le slice orders
	for _, order := range orders {
		if order.ID == id {
			// 4. Si trouvé : encoder et retourner la commande
			json.NewEncoder(w).Encode(order)
			return
		}
	}
	// 4. Si trouvé : encoder et retourner la commande

	// 5. Sinon : retourner une erreur 404
	http.Error(w, "Commande non trouvée", http.StatusNotFound)

}
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Récupérer l'ID depuis les variables de route
	vars := mux.Vars(r)
	id := vars["id"]
	// 2. Parcourir orders avec l'index
	for index, order := range orders {
		if order.ID == id {
			// 3a. Vérifier que le statut n'est pas "picked-up"
			if order.Status == models.StatusPickedUp {
				// 3b. Si picked-up : retourner 400 Bad Request
				http.Error(w, "Impossible de supprimer une commande déjà récupérée", http.StatusBadRequest)
				return
			}
			// 3c. Sinon : supprimer la commande du slice
			orders = append(orders[:index], orders[index+1:]...)
			// 3d. Retourner 204 No Content
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	// 4. Si non trouvée : retourner 404
	http.Error(w, "Commande non trouvée", http.StatusNotFound)
}
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")

	// 2. Récupérer l'ID depuis les variables de route
	vars := mux.Vars(r)
	id := vars["id"]

	// 3. Créer une struct temporaire avec un champ Status
	type StatusUpdate struct {
		Status models.OrderStatus `json:"status"`
	}
	var statusUpdate StatusUpdate
	err := json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// 6. Parcourir orders et trouver la commande
	for index, order := range orders {
		if order.ID == id {
			// 7. Mettre à jour le statut de la commande
			orders[index].Status = statusUpdate.Status
			// 8. Retourner la commande mise à jour en JSON
			json.NewEncoder(w).Encode(orders[index])
			return
		}
	}

	// 9. Si non trouvée : retourner 404
	http.Error(w, "Commande non trouvée", http.StatusNotFound)

}
func main() {
	// TODO 1 : Initialiser le slice drinks avec au moins 5-6 boissons
	drinks = []models.Drink{
		{ID: "DRK-001", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
		{ID: "DRK-002", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
		{ID: "DRK-003", Name: "Latte", Category: "coffee", BasePrice: 4.00},
		{ID: "DRK-004", Name: "Americano", Category: "coffee", BasePrice: 3.00},
		{ID: "DRK-005", Name: "Green Tea", Category: "tea", BasePrice: 2.50},
		{ID: "DRK-006", Name: "Iced Coffee", Category: "cold", BasePrice: 4.50},
	}

	// Exemples : Espresso (2.50€), Cappuccino (3.50€), Latte (4.00€),
	//            Americano (3.00€), Green Tea (2.50€), Iced Coffee (4.50€)

	// TODO 2 : Créer le routeur Mux
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	// TODO 3 : Définir les routes :
	router.HandleFunc("/menu", getMenu).Methods("GET")
	router.HandleFunc("/menu/{id}", getDrink).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")
	router.HandleFunc("/orders", getOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	router.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")
	router.HandleFunc("/orders/{id}/status", UpdateOrderStatus).Methods("PATCH")

	// TODO 4 : (Optionnel) Ajouter une route GET / pour un message de bienvenue
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Bienvenue au service de commande de boissons!"})
	}).Methods("GET")

	// TODO 5 : Afficher un message indiquant que le serveur démarre
	fmt.Println("Démarrage du serveur sur le port 8080...")

	// TODO 6 : Démarrer le serveur sur le port 8080
	http.ListenAndServe(":8080", corsMiddleware(router))
}
