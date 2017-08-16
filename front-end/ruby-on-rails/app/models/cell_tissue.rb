class CellTissue < ApplicationRecord
	belongs_to :cell
	belongs_to :tissue
end