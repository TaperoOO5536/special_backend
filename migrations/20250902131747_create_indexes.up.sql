CREATE INDEX idx_user_events_user_id ON User_Events(User_ID);
CREATE INDEX idx_user_events_event_id ON User_Events(Event_ID);
CREATE INDEX idx_event_pictures_event_id ON Event_Pictures(Event_ID);
CREATE INDEX idx_orders_user_id ON Orders(User_ID);
CREATE INDEX idx_order_items_order_id ON Order_Items(Order_ID);
CREATE INDEX idx_order_items_item_id ON Order_Items(Item_ID);
CREATE INDEX idx_item_pictures_item_id ON Item_Pictures(Item_ID);