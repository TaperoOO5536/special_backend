CREATE INDEX idx_user_ivents_user_id ON User_Ivents(User_ID);
CREATE INDEX idx_user_ivents_ivent_id ON User_Ivents(Ivent_ID);
CREATE INDEX idx_ivent_pictures_ivent_id ON Ivent_Pictures(Ivent_ID);
CREATE INDEX idx_orders_user_id ON Orders(User_ID);
CREATE INDEX idx_order_items_order_id ON Order_Items(Order_ID);
CREATE INDEX idx_order_items_item_id ON Order_Items(Item_ID);
CREATE INDEX idx_item_pictures_item_id ON Item_Pictures(Item_ID);