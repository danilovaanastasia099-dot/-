package repository

import (
	"fmt"
)

type Repository struct {
	services []Service
}

type Service struct {
	ID          int
	Name        string
	Description string
	Type        string  // CPU, GPU, PSU
	TDP         int     // тепловыделение в ваттах
	ImageURL    string  // URL изображения из Minio
	VideoURL    string  // URL видео из Minio
	Status      string  // draft, published, deleted
	Likes       []int   // ID пользователей, поставивших лайк
}

func NewRepository() (*Repository, error) {
	services := []Service{
		{
			ID:          1,
			Name:        "Intel Xeon E5-2680 v4",
			Description: "14-ядерный серверный процессор для стоечных решений. Надёжный выбор для виртуализации и вычислительных узлов.",
			Type:        "CPU",
			TDP:         120,
			ImageURL:    "http://localhost:9000/server-images/cpu-xeon.jpg",
			VideoURL:    "http://localhost:9000/server-videos/cpu-xeon.mp4",
			Status:      "published",
			Likes:       []int{1, 2, 3},
		},
		{
			ID:          2,
			Name:        "NVIDIA Tesla V100",
			Description: "Графический ускоритель для вычислений с TDP 250 Вт. Подходит для ML/AI вычислений и вычислительных шардов.",
			Type:        "GPU",
			TDP:         250,
			ImageURL:    "http://localhost:9000/server-images/gpu-tesla.jpg",
			VideoURL:    "http://localhost:9000/server-videos/gpu-tesla.mp4",
			Status:      "published",
			Likes:       []int{1, 2},
		},
		{
			ID:          3,
			Name:        "Блок питания 800W",
			Description: "Серверный блок питания мощностью 800 Вт с высоким КПД и активным PFC.",
			Type:        "PSU",
			TDP:         800,
			ImageURL:    "http://localhost:9000/server-images/psu-800.jpg",
			VideoURL:    "http://localhost:9000/server-videos/psu-800.mp4",
			Status:      "published",
			Likes:       []int{1},
		},
		{
			ID:          4,
			Name:        "AMD EPYC 7742",
			Description: "64-ядерный серверный процессор с TDP 225 Вт. Отличается высокой плотностью вычислений и производительностью на ядро.",
			Type:        "CPU",
			TDP:         225,
			ImageURL:    "http://localhost:9000/server-images/cpu-epyc.jpg",
			VideoURL:    "http://localhost:9000/server-videos/cpu-epyc.mp4",
			Status:      "draft",
			Likes:       []int{},
		},
		{
			ID:          5,
			Name:        "NVIDIA RTX A6000",
			Description: "Профессиональная видеокарта с высокой производительностью рендеринга и TDP 300 Вт.",
			Type:        "GPU",
			TDP:         300,
			ImageURL:    "http://localhost:9000/server-images/gpu-rtx.jpg",
			VideoURL:    "http://localhost:9000/server-videos/gpu-rtx.mp4",
			Status:      "published",
			Likes:       []int{4, 5},
		},
		{
			ID:          6,
			Name:        "Samsung PM1733 3.84TB",
			Description: "NVMe SSD для серверов, высокопроизводительное хранилище с низкой латентностью. TDP условный 12 Вт.",
			Type:        "STORAGE",
			TDP:         12,
			ImageURL:    "http://localhost:9000/server-images/ssd-pm1733.jpg",
			VideoURL:    "http://localhost:9000/server-videos/ssd-pm1733.mp4",
			Status:      "published",
			Likes:       []int{2, 3},
		},
		{
			ID:          7,
			Name:        "Mellanox ConnectX-6",
			Description: "Сетевая карта 200GbE/100GbE для высокоскоростных кластеров. Низкая задержка, высокая пропускная способность.",
			Type:        "NIC",
			TDP:         25,
			ImageURL:    "http://localhost:9000/server-images/nic-connectx6.jpg",
			VideoURL:    "http://localhost:9000/server-videos/nic-connectx6.mp4",
			Status:      "published",
			Likes:       []int{1},
		},
		{
			ID:          8,
			Name:        "Kingston 128GB DDR4 ECC",
			Description: "Модуль оперативной памяти ECC для серверных систем, 128GB, DDR4. Надёжность и совместимость.",
			Type:        "RAM",
			TDP:         8,
			ImageURL:    "http://localhost:9000/server-images/ram-128gb.jpg",
			VideoURL:    "http://localhost:9000/server-videos/ram-128gb.mp4",
			Status:      "published",
			Likes:       []int{1,2},
		},
	}

	return &Repository{
		services: services,
	}, nil
}

// GetPublishedServices возвращает все опубликованные услуги
func (r *Repository) GetPublishedServices() ([]Service, error) {
	var result []Service
	for _, service := range r.services {
		if service.Status == "published" {
			result = append(result, service)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("нет опубликованных услуг")
	}
	return result, nil
}

// GetServiceByID возвращает услугу по ID
func (r *Repository) GetServiceByID(id int) (Service, error) {
	for _, service := range r.services {
		if service.ID == id && service.Status != "deleted" {
			return service, nil
		}
	}
	return Service{}, fmt.Errorf("услуга не найдена")
}

// GetNextService возвращает следующую услугу после указанного ID
func (r *Repository) GetNextService(id int) (Service, error) {
	for i, service := range r.services {
		if service.ID == id {
			for j := i + 1; j < len(r.services); j++ {
				if r.services[j].Status == "published" {
					return r.services[j], nil
				}
			}
			break
		}
	}
	return Service{}, fmt.Errorf("следующая услуга не найдена")
}

// GetDraftService возвращает черновик
func (r *Repository) GetDraftService() (Service, error) {
	for _, service := range r.services {
		if service.Status == "draft" {
			return service, nil
		}
	}
	return Service{}, fmt.Errorf("черновик не найден")
}

// FilterServicesByTDP фильтрует услуги по TDP
func (r *Repository) FilterServicesByTDP(minTDP int) ([]Service, error) {
	var result []Service
	for _, service := range r.services {
		if service.Status == "published" && service.TDP >= minTDP {
			result = append(result, service)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("услуги не найдены")
	}
	return result, nil
}

// CalculateTotalTDP вычисляет общее тепловыделение
func (r *Repository) CalculateTotalTDP() int {
	total := 0
	for _, service := range r.services {
		if service.Status == "published" {
			total += service.TDP
		}
	}
	return total
}

// ConvertToBTU конвертирует ватты в BTU/час
func ConvertToBTU(watts int) float64 {
	return float64(watts) * 3.412141633
}
