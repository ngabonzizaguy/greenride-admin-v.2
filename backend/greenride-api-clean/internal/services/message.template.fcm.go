package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// FCM消息模板定义 - 系统自带模板
var (
	// 英文FCM模板
	DefaultPassengerOrderAcceptedFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerOrderAccepted,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Driver Accepted",
		Content:     "Vehicle with license plate {{.PlateNumber}} has accepted your ride and is on the way to {{.PickupAddress}}",
		Status:      protocol.StatusActive,
		Description: "Notification when driver accepts passenger's order",
	}

	DefaultPassengerDriverArrivedFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerDriverArrived,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Driver Arrived",
		Content:     "Vehicle with license plate {{.PlateNumber}} has arrived at {{.PickupAddress}}. Please be ready to board.",
		Status:      protocol.StatusActive,
		Description: "Notification when driver arrives at pickup location",
	}

	DefaultPassengerTripStartedFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerTripStarted,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Trip Started",
		Content:     "Your trip has started from {{.PickupAddress}} to {{.DropoffAddress}}. Please fasten your seatbelt.",
		Status:      protocol.StatusActive,
		Description: "Notification when trip starts",
	}

	DefaultPassengerTripEndedFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerTripEnded,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Trip Ended",
		Content:     "You have arrived at {{.DropoffAddress}}. Trip cost: {{.Amount}} {{.Currency}}. Thank you for riding with us.",
		Status:      protocol.StatusActive,
		Description: "Notification when trip ends",
	}

	DefaultPassengerOrderCancelledFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerOrderCancelled,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Ride Cancelled",
		Content:     "Your ride was cancelled by {{.CancelledBy}}. Reason: {{.CancelReason}}",
		Status:      protocol.StatusActive,
		Description: "Notification when ride is cancelled",
	}

	DefaultDriverNewOrderFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverNewOrder,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "New Ride Request",
		Content:     "Passenger {{.PassengerName}} from {{.PickupAddress}} to {{.DropoffAddress}}, distance {{.Distance}}km, estimated {{.Duration}} mins, fare {{.Amount}} {{.Currency}}",
		Status:      protocol.StatusActive,
		Description: "Notification when new ride request is received",
	}

	DefaultDriverTripEndedFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverTripEnded,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Trip Completed",
		Content:     "Trip with passenger {{.PassengerName}} has been completed. Please wait for payment confirmation.",
		Status:      protocol.StatusActive,
		Description: "Notification when trip is completed",
	}

	DefaultDriverPaymentConfirmedFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverPaymentConfirmed,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Payment Confirmed",
		Content:     "Payment of {{.Amount}} {{.Currency}} has been confirmed. Order completed.",
		Status:      protocol.StatusActive,
		Description: "Notification when payment is confirmed",
	}

	DefaultDriverOrderCancelledFcmEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverOrderCancelled,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangEnglish,
		Title:       "Ride Cancelled",
		Content:     "Ride with passenger {{.PassengerName}} was cancelled. Reason: {{.CancelReason}}",
		Status:      protocol.StatusActive,
		Description: "Notification when ride is cancelled",
	}

	// 法语FCM模板
	DefaultPassengerOrderAcceptedFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerOrderAccepted,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Chauffeur accepté",
		Content:     "Le véhicule avec la plaque d'immatriculation {{.PlateNumber}} a accepté votre course et se dirige vers {{.PickupAddress}}",
		Status:      protocol.StatusActive,
		Description: "Notification when driver accepts passenger's order (French)",
	}

	DefaultPassengerDriverArrivedFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerDriverArrived,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Chauffeur arrivé",
		Content:     "Le véhicule avec la plaque d'immatriculation {{.PlateNumber}} est arrivé à {{.PickupAddress}}. Veuillez vous préparer à monter.",
		Status:      protocol.StatusActive,
		Description: "Notification when driver arrives at pickup location (French)",
	}

	DefaultPassengerTripStartedFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerTripStarted,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Trajet commencé",
		Content:     "Votre trajet a commencé de {{.PickupAddress}} à {{.DropoffAddress}}. Veuillez attacher votre ceinture de sécurité.",
		Status:      protocol.StatusActive,
		Description: "Notification when trip starts (French)",
	}

	DefaultPassengerTripEndedFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerTripEnded,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Trajet terminé",
		Content:     "Vous êtes arrivé à {{.DropoffAddress}}. Coût du trajet: {{.Amount}} {{.Currency}}. Merci d'avoir voyagé avec nous.",
		Status:      protocol.StatusActive,
		Description: "Notification when trip ends (French)",
	}

	DefaultPassengerOrderCancelledFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerOrderCancelled,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Course annulée",
		Content:     "Votre course a été annulée par {{.CancelledBy}}. Raison: {{.CancelReason}}",
		Status:      protocol.StatusActive,
		Description: "Notification when ride is cancelled (French)",
	}

	DefaultDriverNewOrderFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverNewOrder,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Nouvelle demande de course",
		Content:     "Passager {{.PassengerName}} de {{.PickupAddress}} à {{.DropoffAddress}}, distance {{.Distance}}km, estimé à {{.Duration}} mins, tarif {{.Amount}} {{.Currency}}",
		Status:      protocol.StatusActive,
		Description: "Notification when new ride request is received (French)",
	}

	DefaultDriverTripEndedFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverTripEnded,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Trajet terminé",
		Content:     "Le trajet avec le passager {{.PassengerName}} est terminé. Veuillez attendre la confirmation du paiement.",
		Status:      protocol.StatusActive,
		Description: "Notification when trip is completed (French)",
	}

	DefaultDriverPaymentConfirmedFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverPaymentConfirmed,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Paiement confirmé",
		Content:     "Le paiement de {{.Amount}} {{.Currency}} a été confirmé. Commande terminée.",
		Status:      protocol.StatusActive,
		Description: "Notification when payment is confirmed (French)",
	}

	DefaultDriverOrderCancelledFcmFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverOrderCancelled,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangFrench,
		Title:       "Course annulée",
		Content:     "La course avec le passager {{.PassengerName}} a été annulée. Raison: {{.CancelReason}}",
		Status:      protocol.StatusActive,
		Description: "Notification when ride is cancelled (French)",
	}

	// 中文FCM模板
	DefaultPassengerOrderAcceptedFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerOrderAccepted,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "司机已接单",
		Content:     "车牌号{{.PlateNumber}}已接单，正在前往{{.PickupAddress}}",
		Status:      protocol.StatusActive,
		Description: "Notification when driver accepts passenger's order (Chinese)",
	}

	DefaultPassengerDriverArrivedFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerDriverArrived,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "司机已到达",
		Content:     "车牌号{{.PlateNumber}}已到达{{.PickupAddress}}，请准备上车",
		Status:      protocol.StatusActive,
		Description: "Notification when driver arrives at pickup location (Chinese)",
	}

	DefaultPassengerTripStartedFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerTripStarted,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "行程开始",
		Content:     "您的行程已开始，从{{.PickupAddress}}前往{{.DropoffAddress}}，请系好安全带",
		Status:      protocol.StatusActive,
		Description: "Notification when trip starts (Chinese)",
	}

	DefaultPassengerTripEndedFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerTripEnded,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "行程结束",
		Content:     "您已到达{{.DropoffAddress}}，行程费用{{.Amount}}{{.Currency}}，感谢您的使用",
		Status:      protocol.StatusActive,
		Description: "Notification when trip ends (Chinese)",
	}

	DefaultPassengerOrderCancelledFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypePassengerOrderCancelled,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "订单已取消",
		Content:     "您的订单已被{{.CancelledBy}}取消，原因：{{.CancelReason}}",
		Status:      protocol.StatusActive,
		Description: "Notification when ride is cancelled (Chinese)",
	}

	DefaultDriverNewOrderFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverNewOrder,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "新订单",
		Content:     "乘客{{.PassengerName}}从{{.PickupAddress}}到{{.DropoffAddress}}，距离{{.Distance}}km，预计{{.Duration}}分钟，金额{{.Amount}}{{.Currency}}",
		Status:      protocol.StatusActive,
		Description: "Notification when new ride request is received (Chinese)",
	}

	DefaultDriverTripEndedFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverTripEnded,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "行程结束",
		Content:     "乘客{{.PassengerName}}的行程已完成，请等待乘客确认支付",
		Status:      protocol.StatusActive,
		Description: "Notification when trip is completed (Chinese)",
	}

	DefaultDriverPaymentConfirmedFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverPaymentConfirmed,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "支付确认",
		Content:     "乘客已确认支付{{.Amount}}{{.Currency}}，订单完成",
		Status:      protocol.StatusActive,
		Description: "Notification when payment is confirmed (Chinese)",
	}

	DefaultDriverOrderCancelledFcmZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeDriverOrderCancelled,
		Channel:     protocol.MsgChannelFcm,
		Language:    protocol.LangChinese,
		Title:       "订单已取消",
		Content:     "乘客{{.PassengerName}}取消了订单，原因：{{.CancelReason}}",
		Status:      protocol.StatusActive,
		Description: "Notification when ride is cancelled (Chinese)",
	}

	// 默认FCM模板集合
	DefaultFcmTemplates = []*models.MessageTemplate{
		// 英文模板
		DefaultPassengerOrderAcceptedFcmEN,
		DefaultPassengerDriverArrivedFcmEN,
		DefaultPassengerTripStartedFcmEN,
		DefaultPassengerTripEndedFcmEN,
		DefaultPassengerOrderCancelledFcmEN,
		DefaultDriverNewOrderFcmEN,
		DefaultDriverTripEndedFcmEN,
		DefaultDriverPaymentConfirmedFcmEN,
		DefaultDriverOrderCancelledFcmEN,

		// 法语模板
		DefaultPassengerOrderAcceptedFcmFR,
		DefaultPassengerDriverArrivedFcmFR,
		DefaultPassengerTripStartedFcmFR,
		DefaultPassengerTripEndedFcmFR,
		DefaultPassengerOrderCancelledFcmFR,
		DefaultDriverNewOrderFcmFR,
		DefaultDriverTripEndedFcmFR,
		DefaultDriverPaymentConfirmedFcmFR,
		DefaultDriverOrderCancelledFcmFR,

		// 中文模板
		DefaultPassengerOrderAcceptedFcmZH,
		DefaultPassengerDriverArrivedFcmZH,
		DefaultPassengerTripStartedFcmZH,
		DefaultPassengerTripEndedFcmZH,
		DefaultPassengerOrderCancelledFcmZH,
		DefaultDriverNewOrderFcmZH,
		DefaultDriverTripEndedFcmZH,
		DefaultDriverPaymentConfirmedFcmZH,
		DefaultDriverOrderCancelledFcmZH,
	}
)
